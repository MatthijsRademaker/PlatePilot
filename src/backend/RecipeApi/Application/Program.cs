using Application;
using Application.DatabaseSeed;
using Application.Endpoints.V1;
using Asp.Versioning;
using Infrastructure;
using Microsoft.EntityFrameworkCore;
using ServiceDefaults;

var builder = WebApplication.CreateBuilder(args);

// Add service defaults & Aspire client integrations.
builder.AddServiceDefaults();
builder.AddApplicationServices();
builder.Services.AddCors();

// Add services to the container.
builder.Services.AddProblemDetails();
var withApiVersioning = builder.Services.AddApiVersioning(options =>
{
    options.ReportApiVersions = true;
    options.AssumeDefaultVersionWhenUnspecified = true;
    options.DefaultApiVersion = new ApiVersion(1, 0);
});

builder.AddDefaultOpenApi(withApiVersioning);
var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    var services = scope.ServiceProvider;
    try
    {
        var context = services.GetRequiredService<RecipeContext>();
        await context.MigrateAsync();

        var seeder = new DatabaseSeeder(context);
        await seeder.SeedAsync();
    }
    catch (Exception ex)
    {
        // Log the error or handle it as needed
        Console.WriteLine($"An error occurred while migrating the database: {ex.Message}");
    }
}

// Configure the HTTP request pipeline.
app.UseExceptionHandler();

app.MapDefaultEndpoints();

app.UseCors(options =>
{
    options.AllowAnyOrigin().AllowAnyHeader().AllowAnyMethod();
});

app.UseStatusCodePages();

app.NewVersionedApi("Recipe").MapRecipeV1();
app.UseDefaultOpenApi();
app.Run();
