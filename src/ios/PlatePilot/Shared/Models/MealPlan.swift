import Foundation

enum MealType: String, CaseIterable, Identifiable, Hashable {
    case breakfast
    case lunch
    case dinner
    case snack

    var id: String { rawValue }

    var label: String {
        rawValue.capitalized
    }

    var systemImage: String {
        switch self {
        case .breakfast:
            return "sunrise"
        case .lunch:
            return "fork.knife"
        case .dinner:
            return "moon.stars"
        case .snack:
            return "takeoutbag.and.cup.and.straw"
        }
    }
}

struct MealSlot: Identifiable, Hashable {
    let id: UUID
    var mealType: MealType
    var recipe: Recipe?
}

struct DayPlan: Identifiable, Hashable {
    let id: UUID
    var date: Date
    var meals: [MealSlot]
}

struct WeekPlan: Identifiable, Hashable {
    let id: UUID
    var startDate: Date
    var endDate: Date
    var days: [DayPlan]
}

extension Date {
    func dayLabel() -> String {
        formatted(.dateTime.weekday(.abbreviated))
    }

    func shortDate() -> String {
        formatted(.dateTime.month(.abbreviated).day())
    }

    func rangeLabel(to endDate: Date) -> String {
        let startMonth = formatted(.dateTime.month(.abbreviated))
        let endMonth = endDate.formatted(.dateTime.month(.abbreviated))
        let startDay = formatted(.dateTime.day())
        let endDay = endDate.formatted(.dateTime.day())
        let year = formatted(.dateTime.year())

        if startMonth == endMonth {
            return "\(startMonth) \(startDay) - \(endDay), \(year)"
        }

        return "\(startMonth) \(startDay) - \(endMonth) \(endDay), \(year)"
    }
}
