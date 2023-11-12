package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var projectTag = "trends-collector"

func CreateVPC(ctx *pulumi.Context, cidr string, resourceName string) (*ec2.Vpc, error) {
	return ec2.NewVpc(ctx, resourceName, &ec2.VpcArgs{
		CidrBlock:          pulumi.String(cidr),
		EnableDnsSupport:   pulumi.Bool(true),
		EnableDnsHostnames: pulumi.Bool(true),
		Tags:               createNameTag(resourceName),
	})
}

func CreateSubnet(
	ctx *pulumi.Context,
	vpc *ec2.Vpc,
	cidr string,
	// availabilityZone string,
	resourceName string,
) (*ec2.Subnet, error) {
	return ec2.NewSubnet(ctx, resourceName, &ec2.SubnetArgs{
		VpcId:     vpc.ID(),
		CidrBlock: pulumi.String(cidr),
		// AvailabilityZone: pulumi.String(availabilityZone),
		Tags: createNameTag(resourceName),
	})
}

func CreateIGW(
	ctx *pulumi.Context, vpc *ec2.Vpc, resourceName string,
) (*ec2.InternetGateway, error) {

	return ec2.NewInternetGateway(ctx, resourceName, &ec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags:  createNameTag(resourceName),
	})
}

func CreatePublicRouteTable(
	ctx *pulumi.Context, vpc *ec2.Vpc, igw *ec2.InternetGateway, resourceName string,
) (*ec2.RouteTable, error) {
	return ec2.NewRouteTable(
		ctx, resourceName, &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: igw.ID(),
				},
			},
			Tags: createNameTag(resourceName),
		},
		pulumi.DependsOn([]pulumi.Resource{vpc, igw}),
	)
}

func CreateRouteTableAssociation(
	ctx *pulumi.Context, routeTable *ec2.RouteTable, subnet *ec2.Subnet, resourceName string,
) (*ec2.RouteTableAssociation, error) {
	return ec2.NewRouteTableAssociation(
		ctx,
		resourceName,
		&ec2.RouteTableAssociationArgs{
			RouteTableId: routeTable.ID(),
			SubnetId:     subnet.ID(),
		},
		pulumi.DependsOn([]pulumi.Resource{routeTable, subnet}),
	)
}

func CreateSecurityGroupForECSTask(
	ctx *pulumi.Context, vpc *ec2.Vpc, resourceName string,
) (*ec2.SecurityGroup, error) {
	return ec2.NewSecurityGroup(
		ctx,
		resourceName,
		&ec2.SecurityGroupArgs{
			VpcId: vpc.ID(),
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					Description: pulumi.String("All outbound traffic"),
					Protocol:    pulumi.String("-1"),
					FromPort:    pulumi.Int(0),
					ToPort:      pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
				},
			},
			Tags: createNameTag(resourceName),
		})
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// VPC //////////////////////////////////////////////////////////////////////////
		resourceName := fmt.Sprintf("%s-vpc", projectTag)
		vpc, err := CreateVPC(ctx, "10.2.0.0/16", resourceName)
		if err != nil {
			return fmt.Errorf("failed create vpc: %v", err)
		}
		ctx.Export(resourceName, vpc.ID())

		// Subnet /////////////////////////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-subnet-app-container", projectTag)
		subnetAppContainer, err := CreateSubnet(ctx, vpc, "10.2.0.0/24", resourceName)
		if err != nil {
			return fmt.Errorf("failed create subnet for App Container: %v", err)
		}
		ctx.Export(resourceName, subnetAppContainer.ID())

		// InternetGateway //////////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-igw", projectTag)
		igw, err := CreateIGW(ctx, vpc, resourceName)
		if err != nil {
			return fmt.Errorf("failed create igw: %v", err)
		}
		ctx.Export(resourceName, igw.ID())

		// ルートテーブル /////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-route-table-public", projectTag)
		publicRouteTable, err := CreatePublicRouteTable(ctx, vpc, igw, resourceName)
		if err != nil {
			return fmt.Errorf("failed create public route table: %v", err)
		}
		ctx.Export(resourceName, publicRouteTable.ID())

		// ルートテーブル 関連付け///////////////////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-route-table-association-app-container", projectTag)
		routeTableAssociationAppContainer, err := CreateRouteTableAssociation(
			ctx, publicRouteTable, subnetAppContainer, resourceName)
		if err != nil {
			return fmt.Errorf("failed create public route association for AppContainer: %v", err)
		}
		ctx.Export(resourceName, routeTableAssociationAppContainer.ID())

		// IAM //////////////////////////////////////////////////////////////
		// ECSタスクロール
		resourceName = fmt.Sprintf("%s-iam-role-for-ecs-task", projectTag)
		ecsTaskRole, err := iam.NewRole(
			ctx,
			resourceName,
			&iam.RoleArgs{
				AssumeRolePolicy: pulumi.String(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Principal": {
							"Service": "ecs-tasks.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}]
				}`),
				InlinePolicies: iam.RoleInlinePolicyArray{
					&iam.RoleInlinePolicyArgs{
						Name: pulumi.String("ecs-task-policy"),
						Policy: pulumi.String(`{
							"Version": "2012-10-17",
							"Statement": [
								{
								   "Effect": "Allow",
								   "Action": [
										"ssmmessages:CreateControlChannel",
										"ssmmessages:CreateDataChannel",
										"ssmmessages:OpenControlChannel",
										"ssmmessages:OpenDataChannel"
								   ],
								  "Resource": "*"
								},
								{
									"Effect": "Allow",
									"Action": [
										"dynamodb:GetItem",
										"dynamodb:Query",
										"dynamodb:Scan",
										"dynamodb:DeleteItem",
										"dynamodb:UpdateItem",
										"dynamodb:PutItem"
									],
								  "Resource": "*"
								}
							]
						}`),
					},
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed create iam role for ecs task: %v", err)
		}
		ctx.Export(resourceName, ecsTaskRole.ID())

		// セキュリティグループ securitygroup ///////////////////////////////////////////////
		resourceName = fmt.Sprintf("%s-sg-app-container", projectTag)
		securityGroupForECSTask, err := CreateSecurityGroupForECSTask(
			ctx, vpc, resourceName)
		if err != nil {
			return fmt.Errorf("failed create security group for ecs task: %v", err)
		}
		ctx.Export(resourceName, securityGroupForECSTask.ID())

		return nil
	})
}

func createNameTag(tag string) pulumi.StringMap {
	return pulumi.StringMap{
		"Name": pulumi.String(tag),
	}
}
