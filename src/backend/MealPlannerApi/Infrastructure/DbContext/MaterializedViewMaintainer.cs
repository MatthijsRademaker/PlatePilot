using Microsoft.EntityFrameworkCore;

namespace Infrastructure;

public interface IMaterializedViewService
{
    Task RefreshRecipeViewAsync(CancellationToken ct = default);
}

public class MaterializedViewService(RecipeContext context) : IMaterializedViewService
{
    public async Task RefreshRecipeViewAsync(CancellationToken ct = default)
    {
        await context.Database.ExecuteSqlRawAsync(
            "REFRESH MATERIALIZED VIEW CONCURRENTLY materialized.recipe_view;",
            ct
        );
    }
}
