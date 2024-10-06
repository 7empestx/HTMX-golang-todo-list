import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as cloudfront from 'aws-cdk-lib/aws-cloudfront';
import * as origins from 'aws-cdk-lib/aws-cloudfront-origins';
import * as acm from 'aws-cdk-lib/aws-certificatemanager';
import * as route53 from 'aws-cdk-lib/aws-route53';
import * as route53Targets from 'aws-cdk-lib/aws-route53-targets';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import * as cdk from 'aws-cdk-lib';

export class CdkStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

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
          cachePolicy: cloudfront.CachePolicy.fromCachePolicyId(
            this, 
            "UseOriginCacheControlHeaders-QueryStrings",
            "4cc15a8a-d715-48a4-82b8-cc0b614638fe"
          ),
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
  }
}
