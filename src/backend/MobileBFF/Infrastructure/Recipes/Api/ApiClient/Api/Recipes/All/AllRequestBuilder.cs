// <auto-generated/>
#pragma warning disable CS0618
using Infrastructure.RecipeApi.Client.Models;
using Microsoft.Kiota.Abstractions.Extensions;
using Microsoft.Kiota.Abstractions.Serialization;
using Microsoft.Kiota.Abstractions;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using System.Threading;
using System;
namespace Infrastructure.RecipeApi.Client.Api.Recipes.All
{
    /// <summary>
    /// Builds and executes requests for operations under \api\recipes\all
    /// </summary>
    [global::System.CodeDom.Compiler.GeneratedCode("Kiota", "1.0.0")]
    public partial class AllRequestBuilder : BaseRequestBuilder
    {
        /// <summary>
        /// Instantiates a new <see cref="global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder"/> and sets the default values.
        /// </summary>
        /// <param name="pathParameters">Path parameters for the request</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public AllRequestBuilder(Dictionary<string, object> pathParameters, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/api/recipes/all?pageIndex={pageIndex}&pageSize={pageSize}{&api%2Dversion*}", pathParameters)
        {
        }
        /// <summary>
        /// Instantiates a new <see cref="global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder"/> and sets the default values.
        /// </summary>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        /// <param name="requestAdapter">The request adapter to use to execute the requests.</param>
        public AllRequestBuilder(string rawUrl, IRequestAdapter requestAdapter) : base(requestAdapter, "{+baseurl}/api/recipes/all?pageIndex={pageIndex}&pageSize={pageSize}{&api%2Dversion*}", rawUrl)
        {
        }
        /// <returns>A List&lt;global::Infrastructure.RecipeApi.Client.Models.RecipeResponse2&gt;</returns>
        /// <param name="cancellationToken">Cancellation token to use when cancelling requests</param>
        /// <param name="requestConfiguration">Configuration for the request such as headers, query parameters, and middleware options.</param>
        /// <exception cref="global::Infrastructure.RecipeApi.Client.Models.ProblemDetails">When receiving a 400 status code</exception>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public async Task<List<global::Infrastructure.RecipeApi.Client.Models.RecipeResponse2>?> GetAsync(Action<RequestConfiguration<global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder.AllRequestBuilderGetQueryParameters>>? requestConfiguration = default, CancellationToken cancellationToken = default)
        {
#nullable restore
#else
        public async Task<List<global::Infrastructure.RecipeApi.Client.Models.RecipeResponse2>> GetAsync(Action<RequestConfiguration<global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder.AllRequestBuilderGetQueryParameters>> requestConfiguration = default, CancellationToken cancellationToken = default)
        {
#endif
            var requestInfo = ToGetRequestInformation(requestConfiguration);
            var errorMapping = new Dictionary<string, ParsableFactory<IParsable>>
            {
                { "400", global::Infrastructure.RecipeApi.Client.Models.ProblemDetails.CreateFromDiscriminatorValue },
            };
            var collectionResult = await RequestAdapter.SendCollectionAsync<global::Infrastructure.RecipeApi.Client.Models.RecipeResponse2>(requestInfo, global::Infrastructure.RecipeApi.Client.Models.RecipeResponse2.CreateFromDiscriminatorValue, errorMapping, cancellationToken).ConfigureAwait(false);
            return collectionResult?.AsList();
        }
        /// <returns>A <see cref="RequestInformation"/></returns>
        /// <param name="requestConfiguration">Configuration for the request such as headers, query parameters, and middleware options.</param>
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
        public RequestInformation ToGetRequestInformation(Action<RequestConfiguration<global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder.AllRequestBuilderGetQueryParameters>>? requestConfiguration = default)
        {
#nullable restore
#else
        public RequestInformation ToGetRequestInformation(Action<RequestConfiguration<global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder.AllRequestBuilderGetQueryParameters>> requestConfiguration = default)
        {
#endif
            var requestInfo = new RequestInformation(Method.GET, UrlTemplate, PathParameters);
            requestInfo.Configure(requestConfiguration);
            requestInfo.Headers.TryAdd("Accept", "application/json");
            return requestInfo;
        }
        /// <summary>
        /// Returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
        /// </summary>
        /// <returns>A <see cref="global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder"/></returns>
        /// <param name="rawUrl">The raw URL to use for the request builder.</param>
        public global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder WithUrl(string rawUrl)
        {
            return new global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder(rawUrl, RequestAdapter);
        }
        [global::System.CodeDom.Compiler.GeneratedCode("Kiota", "1.0.0")]
        #pragma warning disable CS1591
        public partial class AllRequestBuilderGetQueryParameters 
        #pragma warning restore CS1591
        {
#if NETSTANDARD2_1_OR_GREATER || NETCOREAPP3_1_OR_GREATER
#nullable enable
            [QueryParameter("api%2Dversion")]
            public string? ApiVersion { get; set; }
#nullable restore
#else
            [QueryParameter("api%2Dversion")]
            public string ApiVersion { get; set; }
#endif
            [QueryParameter("pageIndex")]
            public int? PageIndex { get; set; }
            [QueryParameter("pageSize")]
            public int? PageSize { get; set; }
        }
        /// <summary>
        /// Configuration for the request such as headers, query parameters, and middleware options.
        /// </summary>
        [Obsolete("This class is deprecated. Please use the generic RequestConfiguration class generated by the generator.")]
        [global::System.CodeDom.Compiler.GeneratedCode("Kiota", "1.0.0")]
        public partial class AllRequestBuilderGetRequestConfiguration : RequestConfiguration<global::Infrastructure.RecipeApi.Client.Api.Recipes.All.AllRequestBuilder.AllRequestBuilderGetQueryParameters>
        {
        }
    }
}
#pragma warning restore CS0618
