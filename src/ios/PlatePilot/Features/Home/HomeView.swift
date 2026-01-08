import SwiftUI

struct HomeView: View {
    var body: some View {
        ScrollView {
            VStack(spacing: 20) {
                HomeHeaderView()

                VStack(spacing: 16) {
                    TodayPlanCard()
                    DailyCalorieTrackerView()
                    RecipeSuggestionsView()
                }
                .padding(.horizontal, 16)
                .padding(.bottom, 24)
            }
        }
        .background(PlatePilotTheme.pageGradient)
    }
}

private struct HomeHeaderView: View {
    @Environment(AuthStore.self) private var authStore
    @State private var showingAccount = false

    private var greeting: String {
        let hour = Calendar.current.component(.hour, from: Date())
        switch hour {
        case 0..<12: return "Good morning"
        case 12..<17: return "Good afternoon"
        default: return "Good evening"
        }
    }

    var body: some View {
        ZStack {
            RoundedRectangle(cornerRadius: 28, style: .continuous)
                .fill(PlatePilotTheme.headerGradient)
                .frame(maxWidth: .infinity)
                .padding(.horizontal, 16)
                .padding(.top, 8)
                .padding(.bottom, 4)

            HStack(alignment: .center) {
                VStack(alignment: .leading, spacing: 6) {
                    Text(greeting)
                        .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                        .foregroundStyle(.white.opacity(0.85))

                    Text("What is cooking?")
                        .font(PlatePilotTheme.titleFont(size: 28))
                        .foregroundStyle(.white)
                }

                Spacer()

                Button {
                    showingAccount = true
                } label: {
                    Image(systemName: "person.crop.circle.fill")
                        .font(.system(size: 28, weight: .semibold))
                        .foregroundStyle(.white.opacity(0.9))
                        .frame(width: 46, height: 46)
                        .plateGlass(cornerRadius: 16, tint: .white.opacity(0.25), interactive: true)
                }
                .buttonStyle(.plain)
            }
            .padding(.horizontal, 32)
            .padding(.vertical, 24)
        }
        .padding(.top, 12)
        .sheet(isPresented: $showingAccount) {
            AccountSheetView { 
                Task { await authStore.logout() }
            }
        }
    }
}

private struct AccountSheetView: View {
    let onSignOut: () -> Void
    @Environment(\.dismiss) private var dismiss

    var body: some View {
        VStack(spacing: 20) {
            Image(systemName: "person.crop.circle.fill")
                .font(.system(size: 44, weight: .semibold))
                .foregroundStyle(PlatePilotTheme.accent)
                .frame(width: 80, height: 80)
                .background(PlatePilotTheme.tintWarm, in: RoundedRectangle(cornerRadius: 24, style: .continuous))

            Text("Signed in")
                .font(PlatePilotTheme.titleFont(size: 20))
                .foregroundStyle(PlatePilotTheme.textPrimary)

            Text("Manage your PlatePilot account session.")
                .font(PlatePilotTheme.bodyFont(size: 14))
                .foregroundStyle(PlatePilotTheme.textSecondary)

            Button("Sign Out") {
                onSignOut()
                dismiss()
            }
            .plateGlassButton(prominent: true)

            Button("Close") {
                dismiss()
            }
            .buttonStyle(.bordered)
        }
        .padding(24)
        .presentationDetents([.medium])
    }
}
