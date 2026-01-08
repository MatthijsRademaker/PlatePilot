import SwiftUI

struct RemoteImageView: View {
    let url: URL?
    let cornerRadius: CGFloat

    init(url: URL?, cornerRadius: CGFloat = PlatePilotMetrics.smallRadius) {
        self.url = url
        self.cornerRadius = cornerRadius
    }

    var body: some View {
        AsyncImage(url: url) { phase in
            switch phase {
            case .empty:
                ZStack {
                    Rectangle().fill(PlatePilotTheme.tintWarm)
                    ProgressView()
                }
            case .success(let image):
                image
                    .resizable()
                    .scaledToFill()
            case .failure:
                ZStack {
                    Rectangle().fill(PlatePilotTheme.tintWarm)
                    Image(systemName: "photo")
                        .font(.system(size: 24, weight: .semibold))
                        .foregroundStyle(.secondary)
                }
            @unknown default:
                Rectangle().fill(PlatePilotTheme.tintWarm)
            }
        }
        .clipShape(RoundedRectangle(cornerRadius: cornerRadius, style: .continuous))
    }
}
