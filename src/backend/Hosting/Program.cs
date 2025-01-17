var builder = DistributedApplication.CreateBuilder(args);

var postgres = builder
    .AddPostgres("postgres")
    .WithImage("ankane/pgvector")
    .WithImageTag("latest")
    .WithLifetime(ContainerLifetime.Session);

var recipeDb = postgres.AddDatabase("recipedb");

var rabbitmq = builder.AddRabbitMQ("messaging").WithManagementPlugin();
;

const string launchProfile = "https";

var recipeApi = builder
    .AddProject<Projects.RecipeApplication>("recipe-api")
    .WithReference(recipeDb)
    .WaitFor(recipeDb)
    .WithReference(rabbitmq);

var mealPlannerApi = builder
    .AddProject<Projects.MealPlannerApplication>("meal-planner-api")
    .WithReference(recipeDb)
    .WaitFor(recipeDb)
    .WithReference(rabbitmq);

// TODO caching layer for recipe api
builder
    .AddProject<Projects.WebApiBFF>("web-api-bff")
    .WithReference(recipeApi)
    .WithReference(mealPlannerApi)
    .WithReference(rabbitmq);

builder.Build().Run();
