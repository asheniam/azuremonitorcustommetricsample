# Azure Monitor custom metrics sample

This is sample code which illustrates how to write and read custom metrics to Azure Monitor.

Before you use this sample, you will need the following:
1) Azure subscription
2) Azure resource which is created in one of the following regions (westu2, westcentralus, northeurope, westeurope, eastus, southcentralus, southeastasia)
3) AAD application to the Azure subscription which has permission to write and read custom metrics

In order to run this sample, you will need to perform the following:
1) Update secret.yml with the AAD application credentials and Azure subscription 
2) Run this executable: ./<binaryname> <resourceId> <resourceRegion>