using Application.Endpoints.V1;
using Asp.Versioning;
using Infrastructure;
using MobileBFF.Infrastructure;
using ServiceDefaults;

var builder = WebApplication.CreateBuilder(args);

// Add service defaults & Aspire client integrations.
builder.AddServiceDefaults();

builder.AddInfrastructure();
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

// Configure the HTTP request pipeline.
app.UseExceptionHandler();

app.MapDefaultEndpoints();
app.MapRecipeApiV1();

app.UseCors(options =>
{
    options.AllowAnyOrigin().AllowAnyHeader().AllowAnyMethod();
});

app.UseStatusCodePages();

// app.MapMapMealPlannerV1();
app.UseDefaultOpenApi();
app.Run();
