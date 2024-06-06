import { Stack, StackProps, RemovalPolicy, CfnOutput } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import * as iam from 'aws-cdk-lib/aws-iam';

export class CdkStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // Define a VPC (Virtual Private Cloud) with non-overlapping CIDR blocks
    const vpc = new ec2.Vpc(this, 'MyVpc', {
      cidr: '10.0.0.0/16',
      maxAzs: 3, // Default is all AZs in the region
      enableDnsHostnames: true,
      enableDnsSupport: true,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'PublicSubnet',
          subnetType: ec2.SubnetType.PUBLIC,
          mapPublicIpOnLaunch: true,
        },
      ],
    });

    // Define a Security Group
    const securityGroup = new ec2.SecurityGroup(this, 'MySecurityGroup', {
      vpc,
      allowAllOutbound: true,
    });

    // Allow SSH access from anywhere
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(22), 'Allow SSH access from anywhere');
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80), 'Allow HTTP access from anywhere');
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(443), 'Allow HTTPS access from anywhere');

    // Define the IAM role for the EC2 instance
    const role = new iam.Role(this, 'InstanceRole', {
      assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com'),
    });

    // Attach policies to the role
    role.addManagedPolicy(iam.ManagedPolicy.fromAwsManagedPolicyName('AmazonS3ReadOnlyAccess'));

    // Define user data script
    const userData = ec2.UserData.forLinux();

    // S3 Bucket for Website Hosting
    const HTMXGOBucket = new s3.Bucket(
      this,
      'HTMXGOBucket',
      {
        bucketName: 'htmx-go',
        removalPolicy: RemovalPolicy.RETAIN,
        versioned: true,
        blockPublicAccess: {
          blockPublicAcls: false,
          blockPublicPolicy: false,
          ignorePublicAcls: false,
          restrictPublicBuckets: false,
        },
        publicReadAccess: true,
      },
    );

    new s3deploy.BucketDeployment(this, `DeployHTMXGO`, {
      sources: [
        s3deploy.Source.asset("../GoApp"),
      ],
      destinationBucket: HTMXGOBucket,
    });

    // Add commands to download and extract the Go application from S3
    userData.addCommands(
      'sudo yum update -y',
      'sudo yum install -y golang tar',
      'mkdir -p /home/ec2-user/GoApp',
      `aws s3 sync s3://${HTMXGOBucket.bucketName} /home/ec2-user/GoApp`,
      'ls -l /home/ec2-user/GoApp',
      'cd /home/ec2-user/GoApp',
      'go build -o myapp .',
      './myapp &'
    );

    // Define an EC2 instance
    const instance = new ec2.Instance(this, 'MyInstance', {
      vpc,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO),
      machineImage: ec2.MachineImage.latestAmazonLinux2(),
      securityGroup,
      keyName: 'MacLaptop',
      userData,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PUBLIC,
      },
      role
    });

    // Output the instance public DNS
    new CfnOutput(this, 'InstancePublicDNS', {
      value: instance.instancePublicDnsName,
    });
  }
}
