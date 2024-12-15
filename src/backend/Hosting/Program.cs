var builder = DistributedApplication.CreateBuilder(args);

var postgres = builder
    .AddPostgres("postgres")
    .WithImage("ankane/pgvector")
    .WithImageTag("latest")
    .WithLifetime(ContainerLifetime.Persistent);

var recipeDb = postgres.AddDatabase("recipedb");

const string launchProfile = "https";

builder
    .AddProject<Projects.Recipe_Application>("application")
    .WithReference(recipeDb)
    .WaitFor(recipeDb);

// TODO - Use this in maui project?
// builder
//     .AddProject<Projects.backend_Web>("webfrontend")
//     .WithExternalHttpEndpoints()
//     .WithReference(apiService)
//     .WaitFor(apiService);

builder.Build().Run();
