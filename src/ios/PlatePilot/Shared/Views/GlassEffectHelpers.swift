import SwiftUI

struct PlateGlassGroup<Content: View>: View {
    let spacing: CGFloat
    let content: () -> Content

    init(spacing: CGFloat = 16, @ViewBuilder content: @escaping () -> Content) {
        self.spacing = spacing
        self.content = content
    }

    var body: some View {
        if #available(iOS 26, *) {
            GlassEffectContainer(spacing: spacing) {
                content()
            }
        } else {
            content()
        }
    }
}

extension View {
    @ViewBuilder
    func plateGlass(
        cornerRadius: CGFloat = PlatePilotMetrics.cardRadius,
        tint: Color = .white,
        interactive: Bool = false
    ) -> some View {
        if #available(iOS 26, *) {
            glassEffect(
                .regular.tint(tint).interactive(interactive),
                in: .rect(cornerRadius: cornerRadius)
            )
        } else {
            background(
                .ultraThinMaterial,
                in: RoundedRectangle(cornerRadius: cornerRadius, style: .continuous)
            )
        }
    }

    @ViewBuilder
    func plateGlassButton(prominent: Bool = false) -> some View {
        if #available(iOS 26, *) {
            if prominent {
                buttonStyle(.glassProminent)
            } else {
                buttonStyle(.glass)
            }
        } else {
            if prominent {
                buttonStyle(.borderedProminent)
            } else {
                buttonStyle(.bordered)
            }
        }
    }
}
