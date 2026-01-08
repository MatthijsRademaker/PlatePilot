import Observation
import SwiftUI

@MainActor
@Observable
final class RouterPath {
    var path: [Route] = []

    func push(_ route: Route) {
        path.append(route)
    }

    func reset() {
        path.removeAll()
    }
}

enum Route: Hashable {
    case recipeDetail(id: UUID)
}

@MainActor
@Observable
final class TabRouter {
    private var routers: [AppTab: RouterPath] = [:]

    func router(for tab: AppTab) -> RouterPath {
        if let router = routers[tab] { return router }
        let router = RouterPath()
        routers[tab] = router
        return router
    }

    func binding(for tab: AppTab) -> Binding<[Route]> {
        let router = router(for: tab)
        return Binding(get: { router.path }, set: { router.path = $0 })
    }
}
