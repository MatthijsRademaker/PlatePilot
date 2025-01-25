using Infrastructure.RecipeApi.Client;
using Microsoft.Kiota.Abstractions.Authentication;
using Microsoft.Kiota.Http.HttpClientLibrary;

namespace MobileBFF.Infrastructure.Recipes.Api;

public class RecipeApiClientFactory(HttpClient httpClient)
{
    public RecipeApi GetClient()
    {
        return new RecipeApi(
            new HttpClientRequestAdapter(
                new AnonymousAuthenticationProvider(),
                httpClient: httpClient
            )
        );
    }
}
