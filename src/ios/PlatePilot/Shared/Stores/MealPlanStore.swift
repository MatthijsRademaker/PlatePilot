import Foundation
import Observation

@MainActor
@Observable
final class MealPlanStore {
    private let apiClient: APIClient

    var currentWeek: WeekPlan
    var isLoading = false
    var isSaving = false
    var errorMessage: String?

    init(apiClient: APIClient = APIClient()) {
        self.apiClient = apiClient
        let startDate = Self.startOfWeek(for: Date())
        self.currentWeek = Self.makeWeekPlan(starting: startDate)
    }

    func loadCurrentWeek() async {
        await loadWeek(for: Date())
    }

    func goToToday() {
        Task { await loadWeek(for: Date()) }
    }

    func shiftWeek(by offset: Int) {
        let calendar = Self.calendar
        guard let newDate = calendar.date(byAdding: .day, value: 7 * offset, to: currentWeek.startDate) else {
            return
        }
        Task { await loadWeek(for: newDate) }
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
        Task { await saveWeek() }
    }

    func setRecipe(_ recipe: Recipe, for slotID: UUID) {
        updateSlot(slotID: slotID, recipe: recipe)
        Task { await saveWeek() }
    }

    func clearRecipe(for slotID: UUID) {
        updateSlot(slotID: slotID, recipe: nil)
        Task { await saveWeek() }
    }

    func fetchSuggestions(amount: Int = 5) async -> [UUID] {
        let recipeIDs = plannedRecipeIDs()
        let payload = MealPlanSuggestRequestDTO(
            dailyConstraints: [],
            alreadySelectedRecipeIds: recipeIDs.map(\.uuidString),
            amount: amount
        )

        do {
            let response = try await apiClient.suggestMealPlan(payload: payload)
            let ids = response.recipeIds ?? []
            return ids.compactMap { UUID(uuidString: $0) }
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to fetch meal plan suggestions."
            return []
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

    private func loadWeek(for date: Date) async {
        let startDate = Self.startOfWeek(for: date)
        currentWeek = Self.makeWeekPlan(starting: startDate)
        isLoading = true
        errorMessage = nil

        do {
            let dto = try await apiClient.fetchWeekPlan(startDate: Self.formatDate(startDate))
            currentWeek = Self.weekPlan(from: dto, fallback: currentWeek)
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to load meal plan."
        }

        isLoading = false
    }

    private func saveWeek() async {
        isSaving = true
        errorMessage = nil

        let payload = MealPlanWeekSaveRequestDTO(
            startDate: Self.formatDate(currentWeek.startDate),
            endDate: Self.formatDate(currentWeek.endDate),
            days: currentWeek.days.map { day in
                MealPlanDayInputDTO(
                    date: Self.formatDate(day.date),
                    meals: day.meals.map { slot in
                        MealPlanSlotInputDTO(
                            mealType: slot.mealType.rawValue,
                            recipeId: slot.recipe?.id.uuidString
                        )
                    }
                )
            }
        )

        do {
            let dto = try await apiClient.saveWeekPlan(payload: payload)
            currentWeek = Self.weekPlan(from: dto, fallback: currentWeek)
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to save meal plan."
        }

        isSaving = false
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

    private func plannedRecipeIDs() -> [UUID] {
        var ids: [UUID] = []
        for day in currentWeek.days {
            for meal in day.meals {
                if let recipeID = meal.recipe?.id {
                    ids.append(recipeID)
                }
            }
        }
        return ids
    }
}

private extension MealPlanStore {
    static let calendar: Calendar = {
        var calendar = Calendar(identifier: .gregorian)
        calendar.firstWeekday = 2
        calendar.timeZone = .current
        return calendar
    }()

    static let dateFormatter: DateFormatter = {
        let formatter = DateFormatter()
        formatter.calendar = calendar
        formatter.locale = Locale(identifier: "en_US_POSIX")
        formatter.timeZone = calendar.timeZone
        formatter.dateFormat = "yyyy-MM-dd"
        return formatter
    }()

    static let defaultMealTypes: [MealType] = [.breakfast, .lunch, .dinner]

    static func startOfWeek(for date: Date) -> Date {
        let calendar = Self.calendar
        let components = calendar.dateComponents([.yearForWeekOfYear, .weekOfYear], from: date)
        return calendar.date(from: components) ?? date
    }

    static func makeWeekPlan(starting startDate: Date) -> WeekPlan {
        let calendar = Self.calendar
        let endDate = calendar.date(byAdding: .day, value: 6, to: startDate) ?? startDate

        var days: [DayPlan] = []
        for offset in 0..<7 {
            guard let date = calendar.date(byAdding: .day, value: offset, to: startDate) else { continue }
            let meals = defaultMealTypes.map { mealType in
                MealSlot(id: UUID(), mealType: mealType, recipe: nil)
            }
            days.append(DayPlan(id: UUID(), date: date, meals: meals))
        }

        return WeekPlan(id: UUID(), startDate: startDate, endDate: endDate, days: days)
    }

    static func formatDate(_ date: Date) -> String {
        dateFormatter.string(from: date)
    }

    static func parseDate(_ value: String?) -> Date? {
        guard let value, !value.isEmpty else { return nil }
        return dateFormatter.date(from: value)
    }

    static func weekPlan(from dto: MealPlanWeekDTO, fallback: WeekPlan?) -> WeekPlan {
        let startDate = parseDate(dto.startDate) ?? startOfWeek(for: Date())
        let endDate = parseDate(dto.endDate) ?? calendar.date(byAdding: .day, value: 6, to: startDate) ?? startDate

        let fallbackRecipes = fallback.flatMap(recipeLookup) ?? [:]
        let dayDTOs = dto.days ?? []

        let days = dayDTOs.compactMap { dayDTO -> DayPlan? in
            guard let date = parseDate(dayDTO.date) else { return nil }
            let slots = dayDTO.meals ?? []
            var slotMap: [MealType: Recipe?] = [:]

            for slot in slots {
                guard let rawMealType = slot.mealType,
                      let mealType = MealType(rawValue: rawMealType) else {
                    continue
                }

                if let recipeDTO = slot.recipe,
                   let recipeID = UUID(uuidString: recipeDTO.id ?? "") {
                    if let existing = fallbackRecipes[recipeID] {
                        slotMap[mealType] = existing
                    } else {
                        let name = recipeDTO.name ?? "Untitled Recipe"
                        slotMap[mealType] = Recipe.summary(id: recipeID, name: name, description: recipeDTO.description)
                    }
                } else {
                    slotMap[mealType] = nil
                }
            }

            let meals = defaultMealTypes.map { mealType in
                MealSlot(id: UUID(), mealType: mealType, recipe: slotMap[mealType] ?? nil)
            }

            return DayPlan(id: UUID(), date: date, meals: meals)
        }

        let sortedDays = days.sorted { $0.date < $1.date }

        if sortedDays.isEmpty {
            return makeWeekPlan(starting: startDate)
        }

        return WeekPlan(id: UUID(), startDate: startDate, endDate: endDate, days: sortedDays)
    }

    static func recipeLookup(from plan: WeekPlan) -> [UUID: Recipe] {
        var lookup: [UUID: Recipe] = [:]
        for day in plan.days {
            for meal in day.meals {
                if let recipe = meal.recipe {
                    lookup[recipe.id] = recipe
                }
            }
        }
        return lookup
    }
}
