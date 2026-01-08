import Foundation
import Observation

@MainActor
@Observable
final class MealPlanStore {
    var currentWeek: WeekPlan

    init(currentWeek: WeekPlan = MockData.weekPlan(starting: Date(), recipes: MockData.recipes)) {
        self.currentWeek = currentWeek
    }

    func goToToday() {
        currentWeek = MockData.weekPlan(starting: Date(), recipes: MockData.recipes)
    }

    func shiftWeek(by offset: Int) {
        let calendar = Calendar.current
        guard let newDate = calendar.date(byAdding: .day, value: 7 * offset, to: currentWeek.startDate) else { return }
        currentWeek = MockData.weekPlan(starting: newDate, recipes: MockData.recipes)
    }

    func clearWeek() {
        currentWeek.days = currentWeek.days.map { day in
            var updated = day
            updated.meals = day.meals.map { slot in
                var updatedSlot = slot
                updatedSlot.recipe = nil
                return updatedSlot
            }
            return updated
        }
    }

    func setRecipe(_ recipe: Recipe, for slotID: UUID) {
        updateSlot(slotID: slotID, recipe: recipe)
    }

    func clearRecipe(for slotID: UUID) {
        updateSlot(slotID: slotID, recipe: nil)
    }

    private func updateSlot(slotID: UUID, recipe: Recipe?) {
        currentWeek.days = currentWeek.days.map { day in
            var updatedDay = day
            updatedDay.meals = day.meals.map { slot in
                guard slot.id == slotID else { return slot }
                var updatedSlot = slot
                updatedSlot.recipe = recipe
                return updatedSlot
            }
            return updatedDay
        }
    }

    func featuredMeal() -> MealSlot? {
        let priority: [MealType] = [.dinner, .lunch, .breakfast]
        for day in currentWeek.days where Calendar.current.isDateInToday(day.date) {
            for mealType in priority {
                if let slot = day.meals.first(where: { $0.mealType == mealType && $0.recipe != nil }) {
                    return slot
                }
            }
        }
        return nil
    }
}
