var builder = DistributedApplication.CreateBuilder(args);

var postgres = builder
    .AddPostgres("postgres")
    .WithImage("ankane/pgvector")
    .WithImageTag("latest")
    .WithLifetime(ContainerLifetime.Session)
    .WithPgAdmin(
        (builder) =>
        {
            builder.WithHostPort(5050);
        }
    );

var recipeDb = postgres.AddDatabase("recipedb");

var mssqlInstance = builder
    .AddContainer("mssql", "mcr.microsoft.com/mssql/server:2022-latest", "")
    .WithEnvironment("ACCEPT_EULA", "Y")
    .WithEnvironment("MSSQL_SA_PASSWORD", "temporarily-secure-password-!123");

var serviceBusInstance = builder
    .AddContainer("servicebus", "mcr.microsoft.com/azure-messaging/servicebus-emulator")
    .WithEnvironment("ACCEPT_EULA", "Y")
    .WithEnvironment("SQL_SERVER", mssqlInstance.Resource.Name)
    .WithEnvironment("MSSQL_SA_PASSWORD", "temporarily-secure-password-!123")
    .WithBindMount(
        "servicebus.emulator.config.json",
        "/ServiceBus_Emulator/ConfigFiles/Config.json"
    )
    .WithEndpoint(
        "servicebus",
        (endpoint) =>
        {
            endpoint.Name = "servicebus";
            endpoint.Port = 5672;
            endpoint.TargetPort = 5672;
        }
    )
    .WaitFor(mssqlInstance);

var serviceBus = builder.AddConnectionString("messaging");

const string launchProfile = "https";

var recipeApi = builder
    .AddProject<Projects.RecipeApplication>("recipe-api")
    .WithReference(recipeDb)
    .WaitFor(recipeDb)
    .WaitFor(serviceBusInstance)
    .WithReference(serviceBus);

var mealPlannerApi = builder
    .AddProject<Projects.MealPlannerApplication>("meal-planner-api")
    .WithReference(recipeDb)
    .WaitFor(recipeDb)
    .WaitFor(serviceBusInstance)
    .WithReference(serviceBus);

// TODO caching layer for recipe api
builder
    .AddProject<Projects.WebApiBffApplication>("web-api-bff")
    .WithReference(recipeApi)
    .WithReference(mealPlannerApi)
    .WaitFor(serviceBusInstance)
    .WithReference(serviceBus);

builder.Build().Run();
