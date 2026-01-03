Feature: AI-powered recipe suggestions in meal planning
  As a user planning my weekly meals
  I want to receive intelligent recipe suggestions
  So that I can easily discover diverse recipes that fit my preferences

  Scenario: User requests recipe suggestions for empty meal slot
    Given I am on the meal planning page
    And I have an empty meal slot
    When I click "Get Suggestions" on the meal slot
    Then I see a list of suggested recipes from the AI/vector-based recommendation system
    And the suggestions exclude recipes already planned for the week

  Scenario: User accepts a suggested recipe
    Given I am viewing recipe suggestions for a meal slot
    When I select a suggested recipe
    Then the recipe is assigned to the meal slot
    And the suggestions dialog closes

Technical Notes:
- Use existing mealplanApi.suggestRecipes() endpoint in frontend
- Pass excludeRecipeIds from plannedRecipeIds computed property
- Add "Suggest" button to MealSlotCard component
- Create SuggestionsDialog component for displaying recommendations
- The backend vector search is already implemented

Acceptance Criteria:
- [ ] MealSlotCard has a "Suggest" icon button
- [ ] Clicking suggest opens a dialog with recipe recommendations
- [ ] Suggestions exclude already-planned recipes for the week
- [ ] User can select a suggestion to fill the slot
- [ ] Loading state shown while fetching suggestions
