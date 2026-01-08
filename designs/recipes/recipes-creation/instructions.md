
# PlatePilot — Recipe Creation Screen (SwiftUI) Implementation Spec
Version: 1.0  
Scope: “New Recipe” creation screen matching warm liquid-glass mockup style.

---

## 0) Design Goal

Create a **premium, iOS-native** “Recipe Creation” screen that feels:
- calm and guided (not form-like),
- warm and food-forward,
- “liquid glass” in depth and material,
- highly legible with strong hierarchy.

Keep all functional components from the existing mockup:
- Recipe title input + photo action
- Description input
- Prep/Cook time pickers (combined row)
- Ingredients list + add
- Instructions list + add + delete
- Guided Cooking Mode toggle
- Dietary tags (Vegetarian / Vegan / Gluten-Free)
- Primary “Save Recipe” button

---

## 1) Visual System (Tokens)

### 1.1 Color Palette (Light)
Use warm amber/orange accent without looking “cheap”.

**Background Gradient**
- Top: `#F28C38` (Warm Amber)
- Mid: `#F7B267` (Apricot)
- Bottom: `#FFF2E8` (Cream)

**Text**
- Primary: near-black warm gray `#1D1B18`
- Secondary: `#6E6259`
- Placeholder: `#9A8C82`

**Accent**
- Accent: `#E0701F`
- Accent glow: `#FFB36B` (used sparingly)

### 1.2 Materials (Liquid Glass)
Use SwiftUI materials + subtle overlays:
- Primary surfaces: `.ultraThinMaterial` OR `.thinMaterial` depending on contrast
- Add a soft white overlay (5–10% opacity)
- Add an inner highlight stroke (white @ ~25% opacity)
- Add a soft shadow (warm, low spread)

### 1.3 Corner Radius + Spacing
- Surface radius: 18–22
- Button radius: 24–28 (pill)
- Base horizontal padding: 20
- Section vertical spacing: 18–22
- Internal row padding: 14–16

---

## 2) Layout Structure

### 2.1 Root
- `ZStack` background gradient + subtle grain overlay
- `ScrollView` content
- Sticky bottom action area (safe area aware) for Save button (optional, recommended)

Structure:
- Header (title + microcopy)
- Hero card (title + camera)
- Meta card (description + time row)
- Ingredients section
- Instructions section
- Guided cooking toggle row
- Dietary tags row
- Save CTA (bottom)

---

## 3) Components (SwiftUI Specs)

### 3.1 Background
**Gradient layer**
- `LinearGradient` top→bottom using palette above

**Grain/Noise overlay**
- Option A: use a subtle noise PNG (recommended)
- Option B: procedural noise using Canvas (advanced)
- Opacity: 0.06–0.10
- Blend: `.overlay` or `.softLight`

### 3.2 Header
**Title**
- Text: “New Recipe”
- Weight: semibold/bold
- Size: 30–34
- Tracking: slightly tight (-0.2)

**Subtitle**
- “Let’s build something delicious”
- Size: 15–16
- Opacity: 0.75

Placement:
- Top aligned inside scroll, with extra breathing room
- Ensure not cramped under notch

### 3.3 Glass Surface (Reusable)
Create a reusable `GlassCard`:
- Background: `.ultraThinMaterial`
- Overlay: `Color.white.opacity(0.08)`
- Border stroke: `LinearGradient(white opacity 0.35 → 0.08)` 1px
- Shadow:
  - warm shadow: `Color.orange.opacity(0.18)` y: 10 blur: 25
  - neutral shadow: `Color.black.opacity(0.08)` y: 6 blur: 18
- Corner radius: 20

### 3.4 Recipe Title Row + Photo Action
Inside GlassCard:
- Title `TextField("Recipe Title", text: ...)`
  - large: 20–22
  - semibold
- Trailing photo button:
  - circular (44x44)
  - glass background + accent icon
  - tap → photo picker

Photo button visuals:
- Fill: `Color.white.opacity(0.14)` over material
- Stroke: white @ 0.25
- Icon: SF Symbol `camera.fill` in accent color

### 3.5 Description
- Use `TextEditor` (or custom multiline field)
- Placeholder overlay
- Min height: ~84
- On focus: card subtly “lifts” (see motion)

### 3.6 Time Row (Prep + Cook)
Single horizontal pill row inside the meta card:
- Left segment: Prep time
- Right segment: Cook time
- Each segment:
  - icon: `clock`
  - label + value (e.g., “Prep 15 min”)
- Divider between segments: 1px white @ 0.2

Interaction:
- tap on segment → `bottom sheet` or `sheet` with wheel picker
- Keep this simple: `sheet` with `DateComponents` minutes picker

### 3.7 Ingredients Section
Header row:
- “Ingredients”
- trailing `+` circular button

List:
- Each ingredient row is its own small GlassRow:
  - bullet dot or tiny icon
  - ingredient text
- Swipe actions:
  - swipe trailing → delete
- Tap row → edit inline (optional v1)
- Add:
  - `sheet` with:
    - text field
    - optional quantity/unit
    - save

### 3.8 Instructions Section
Header row:
- “Instructions”
- trailing `+` button (same style)

List:
- Each step row:
  - number badge (1,2,3…)
  - step text
  - trailing optional delete icon OR swipe-to-delete

Number badge:
- rounded rect or circle
- fill: white opacity 0.18
- text: semibold

### 3.9 Guided Cooking Mode Toggle
Row in its own soft glass capsule:
- Label: “Guided Cooking Mode”
- Toggle style:
  - tint: accent
  - consider custom toggle visuals if needed later

Add a short subcaption optionally:
- “Keeps steps on-screen while you cook” (secondary text)

### 3.10 Dietary Tags (Pill Toggles)
Replace switches with pill toggles:
- Vegetarian, Vegan, Gluten-Free
- Each is a `Button` that toggles state
- Inactive:
  - glass fill (white @ 0.10)
  - border white @ 0.20
- Active:
  - gradient fill (accent)
  - subtle glow shadow

Accessibility:
- Provide `.accessibilityLabel("Vegetarian")`
- `.accessibilityValue("On/Off")`

### 3.11 Save Button (Primary CTA)
Large pill button:
- Height: 56–60
- Gradient fill: accent gradient (amber → deeper orange)
- Text: white, semibold, 18
- Subtle gloss highlight (top overlay)
- Press state:
  - scale 0.98
  - darken slightly
- Disabled state:
  - reduce opacity to 0.55
  - no glow

Placement:
- Prefer sticky bottom container with safe-area padding
- Add background blur behind CTA region for readability

---

## 4) Motion + Microinteractions

### 4.1 Screen Entrance
On appear:
- Header fades + moves down slightly (10–14 px)
- Cards stagger:
  - delay per section: 40–60ms
  - duration: 360–520ms
  - easing: `.smooth` (iOS 17) or `.easeOut`

### 4.2 Focus Lift
When any field is focused:
- parent card:
  - shadow increases slightly
  - y offset -2
  - border opacity increases

### 4.3 Add Buttons
`+` button:
- on press: quick spring
- on add success: optional tiny confetti sparkle (very subtle, optional)

### 4.4 Haptics
Use `UIImpactFeedbackGenerator`:
- light: tag toggles
- medium: Save
- rigid: delete confirmation (optional)

---

## 5) State + Data Model (Suggested)

### 5.1 View Model (RecipeDraft)
Fields:
- `title: String`
- `description: String`
- `prepMinutes: Int`
- `cookMinutes: Int`
- `ingredients: [IngredientDraft]`
- `steps: [String]`
- `guidedMode: Bool`
- `tags: Set<DietTag>`
- `photo: UIImage?` or `PhotosPickerItem?`

### 5.2 Validation
- Save disabled unless:
  - title not empty
  - at least 1 ingredient
  - at least 1 instruction step

---

## 6) Accessibility + Dynamic Type

- Use `DynamicTypeSize` friendly layouts:
  - allow list rows to wrap
- Ensure contrast on glass:
  - text must remain readable (avoid too-light text)
- Every icon button has:
  - accessibility label
  - minimum hit area 44x44
- Support Reduce Transparency:
  - if enabled, fall back to solid light fills instead of heavy blur

---

## 7) Implementation Checklist

### Visual
- [ ] Warm gradient background
- [ ] Grain overlay
- [ ] Glass cards with consistent radius/shadows
- [ ] Consistent `+` icon buttons
- [ ] Pill dietary toggles
- [ ] Premium Save CTA (sticky)

### Functional
- [ ] Photo picker for recipe image
- [ ] Prep/cook pickers
- [ ] Add/delete ingredients
- [ ] Add/delete instruction steps
- [ ] Guided mode toggle
- [ ] Tag toggles
- [ ] Save validation

### Polish
- [ ] Entrance animation
- [ ] Focus lift
- [ ] Haptics
- [ ] Accessibility pass

---

## 8) Notes for “PlatePilot Vision” Alignment

This screen is intentionally designed to later support:
- AI agent hints (“Try adding a protein?”)
- Auto-suggested ingredients (autocomplete)
- “Guided Cooking Mode” step-by-step full screen
- Persona-themed skins per agent

Keep this screen neutral-warm and premium; agent personality can be layered via microcopy and small accents without redesigning the base components.

---
