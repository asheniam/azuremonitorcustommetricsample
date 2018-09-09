# Azure Monitor custom metrics sample
<pre>
This is sample code which illustrates how to write and read custom metrics to Azure Monitor.

Before you use this sample, you will need the following:
1) Azure subscription
2) Azure resource which is created in one of the following regions (westu2, westcentralus, northeurope, westeurope, eastus, southcentralus, southeastasia)
3) AAD application to the Azure subscription which has permission to write and read custom metrics (Microsoft.Insights/metrics/*)

In order to run this sample, you will need to perform the following:
1) Compile the source code: go build -o &lt;binaryname&gt; .
2) Update secret.yml with the AAD application credentials and Azure subscription 
3) Run the executable: ./&lt;binaryname&gt; &lt;resourceId&gt; &lt;resourceRegion&gt;
</pre>
