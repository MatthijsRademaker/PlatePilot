using Domain;
using MediatR;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

namespace Application.Endpoints.V1;

public static class MealPlannerEndpoint
{
    public static IEndpointRouteBuilder MapMealPlannerV1(this IEndpointRouteBuilder endpoints)
    {
        // endpoints
        //     .MapGroup("v1")
        //     .MapPost(
        //         "/plan-meal",
        //         (IMediator mediator, PlanMealRequest request) =>
        //         {
        //             return mediator.Send(request);
        //         }
        //     );
        return endpoints;
    }
}
