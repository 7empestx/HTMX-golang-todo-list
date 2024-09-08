import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53Targets from 'aws-cdk-lib/aws-route53-targets';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as cdk from 'aws-cdk-lib';

export class CdkStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    /*
    const hostedZoneDomainName = "grantstarkman.com";

    // Must hardcode the domain name of the Elastic Beanstalk environment since its setup manually
    const elasticBeanstalkEnvironmentDomain = process.env.ELASTIC_BEANSTALK_ENVIRONMENT_DOMAIN || "";
    const htmxGoDomainName = "gohtmxtodo.grantstarkman.com";

    let hostedZone: route53.IHostedZone;
    hostedZone = route53.HostedZone.fromLookup(this, `HostedZone`, {
      domainName: hostedZoneDomainName,
    });

    const HTMXGoCloudfrontSiteCertificate = new acm.Certificate(
      this,
      "HTMXGoCloudfrontSiteCertificate",
      {
        domainName: htmxGoDomainName,
        certificateName: htmxGoDomainName,
        validation: acm.CertificateValidation.fromDns(hostedZone),
      },
    );

    // Cloudfront Distribution for HTMX Go Site 
    const htmxGoDistribution = new cloudfront.Distribution(
      this,
      "HTMXGoDistribution",
      {
        comment: `CloudFront distribution for HTMX Go website`,
        defaultBehavior: {
          origin: new origins.HttpOrigin(elasticBeanstalkEnvironmentDomain, {
            protocolPolicy: cloudfront.OriginProtocolPolicy.HTTP_ONLY,
            httpPort: 80,
          }),
          viewerProtocolPolicy:
            cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
          allowedMethods: cloudfront.AllowedMethods.ALLOW_ALL,
          cachePolicy: cloudfront.CachePolicy.CACHING_DISABLED,
        },
        errorResponses: [
          {
            httpStatus: 403,
            responseHttpStatus: 200,
            responsePagePath: "/index.html",
          },
          {
            httpStatus: 404,
            responseHttpStatus: 200,
            responsePagePath: "/index.html",
          },
        ],
        domainNames: [htmxGoDomainName],
        certificate: HTMXGoCloudfrontSiteCertificate,
      },
    );

    // Route 53 Records for Cloudfront Distribution Frontend
    new route53.ARecord(this, `HTMXGoCloudFrontARecord`, {
      zone: hostedZone,
      recordName: htmxGoDomainName,
      target: route53.RecordTarget.fromAlias(
        new route53Targets.CloudFrontTarget(htmxGoDistribution)
      ),
    });

    new cognito.UserPool(this, 'HTMXGoUserPool', {
      userPoolName: 'HTMXGoUserPool',
      selfSignUpEnabled: true,
      signInAliases: {
        email: true,
      },
      autoVerify: {
        email: true,
      },
      passwordPolicy: {
        minLength: 8,
        requireDigits: false,
        requireLowercase: false,
        requireSymbols: false,
        requireUppercase: false,
      },
      accountRecovery: cognito.AccountRecovery.EMAIL_ONLY,
      userVerification: {
        emailStyle: cognito.VerificationEmailStyle.CODE,
      },
    });
    */

    // Create a VPC
    const vpc = new ec2.Vpc(this, 'GoAppVPC', {
      ipAddresses: ec2.IpAddresses.cidr('172.31.0.0/16'),
      maxAzs: 2,
      enableDnsSupport: true,
      enableDnsHostnames: true,
      subnetConfiguration: [
        {
          cidrMask: 24,
          name: 'Public',
          subnetType: ec2.SubnetType.PUBLIC,
        } 
      ]
    });

    // Create a security group
    const securityGroup = new ec2.SecurityGroup(this, 'GoAppSecurityGroup', {
      vpc,
      description: 'Allow SSH (TCP port 22) and HTTP (TCP port 80) in',
      allowAllOutbound: true
    });
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(22), 'Allow SSH Access')
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(8080), 'Allow HTTP Access')

    // Create a role for the EC2 instance
    const role = new iam.Role(this, 'GoAppInstanceRole', {
      assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com')
    });

    role.addToPolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: [
        'ec2:AuthorizeSecurityGroupIngress',
        'ec2:RevokeSecurityGroupIngress',
        'ec2:AuthorizeSecurityGroupEgress',
        'ec2:RevokeSecurityGroupEgress',
        'ec2:ModifyVpcAttribute',
        'ec2:DescribeSecurityGroups',
        'ec2:DescribeVpcs',
      ],
      resources: ['*'],
    }));

    // Create the EC2 instance
    const ec2Instance = new ec2.Instance(this, 'GoAppInstance', {
      vpc,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.NANO),
      machineImage: new ec2.AmazonLinuxImage({ generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2 }),
      securityGroup: securityGroup,
      role: role,
      keyPair: ec2.KeyPair.fromKeyPairName(this, 'GoAppKeyPair', 'MacLaptop'),
    });

    // Add user data script to install Go and run the app
    const userDataScript = ec2.UserData.forLinux();
    userDataScript.addCommands(
      'yum update -y',
      'yum install -y golang git',
      'mkdir -p /home/ec2-user/go/src/github.com/7empestx',
      'cd /home/ec2-user/go/src/github.com/7empestx',
      'git clone https://github.com/7empestx/HTMX-golang-todo-list.git',
      'cd HTMX-golang-todo-list',
      'go mod tidy',
      'go build -o app cmd/server/main.go',
      'nohup ./app > app.log 2>&1 &',
      'echo "Go app is running on port 8080. Check app.log for details."'
    );
    ec2Instance.addUserData(userDataScript.render());

    // Output the public IP of the EC2 instance
    new cdk.CfnOutput(this, 'InstancePublicIp', {
      value: ec2Instance.instancePublicIp,
      description: 'Public IP address of the EC2 instance',
    });

    // Output the DNS name of the EC2 instance
    new cdk.CfnOutput(this, 'InstanceDnsName', {
      value: ec2Instance.instancePublicDnsName,
      description: 'Dns Name address of the EC2 instance',
    });

  }
}
