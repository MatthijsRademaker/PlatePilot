# Glass Hub Navigation – Apple Liquid Glass Inspired Design

This document describes a radial navigation system inspired by Apple’s “liquid glass” aesthetic, designed around a central hub with orbiting options. It combines glassmorphism, spring-based motion, and contextual navigation modes.

---

## 1. Core Component: “Glass Hub” Radial Navigation

### Layout
- **Center Button (Hub)**
  - Circular glass puck
  - Diameter: 56–64pt
  - Persistent across the app
- **Orbit Items (ö)**
  - 3–6 smaller circular glass bubbles
  - Diameter: 40–48pt
  - Arranged on an arc or full circle around the hub
- **Dynamic Connector (`--`)**
  - Replaced by a liquid, gel-like connector
  - Stretches from hub to the currently hovered or selected bubble
  - Morphs dynamically based on interaction

---

## 2. Visual Style – Liquid Glass / Glassmorphism

### Material Rules
Apply consistently across hub, orbit items, and connector:

- **Background Blur**
  - Strong blur with subtle tint
  - Avoid pure transparency
- **Specular Highlight**
  - Soft highlight gradient near top-left
- **Inner Stroke**
  - 1px white stroke
  - Opacity: 10–18%
- **Outer Shadow**
  - Very soft shadow
  - Low opacity, large blur radius

### Bubble Construction
- Fill: Translucent material (blur + tint)
- Stroke: White at 12–18% opacity
- Shadow: Black at 10–14% opacity
- Highlight: White-to-transparent gradient near top-left

---

## 3. Interaction States

### Idle
- Hub has subtle “breathing” animation
- Scale oscillates from 1.00 → 1.02 over ~2 seconds

### Press Hub
- Orbit bubbles emerge with spring animation
- Slight overshoot before settling

### Hover / Drag (if supported)
- Nearest bubble magnifies slightly
- Connector snaps fluidly to that bubble

### Select Bubble
- Selected bubble locks into place
- Hub glow changes to match selected section’s accent color

### Double Press Hub
- Returns user to Home
- Orbit bubbles collapse back into hub

---

## 4. Motion Design (Apple-like Feel)

- All animations use spring physics with damping
- Orbit bubbles animate with:
  - Scale: 0.7 → 1.0
  - Opacity: 0 → 1
  - Position: Hub center → orbit position
- Connector animation:
  - Morphing blob, not a straight line
  - Thickens slightly on selection
  - Subtle wobble before settling

---

## 5. Liquid Connector (“Dynamic Dashes”)

The `--` connector is implemented as a **glass gel bridge**:

- Rounded, capsule-like shape
- Stretches between hub and selected bubble
- Includes:
  - Faint inner highlight
  - Soft shadow
  - Subtle gradient shift to suggest refraction

### Behavior
- Slides and reshapes when switching between options
- Brief bulge near selected bubble on confirmation
- Smooth morphing between states

---

## 6. Navigation Architecture – Section-Based Modes

The app is divided into 3–4 high-level modes:

### Example Sections
- Home / Overview
- Calories
- Meal Plan
- Recipes
- (Optional) Insights / Progress

### Behavior
- The app operates in **modes**
- The central hub is the mode switcher
- Each mode has its own contextual navigation

---

## 7. Contextual Bottom Navigation Bar

Inspired by Apple Notes:

- **Persistent Bottom Bar**
  - Remains visible across all sections
- **Contextual Items**
  - Left and right actions change per mode
  - Center hub remains constant

### Example Layout
- Left: Primary action (e.g., “Log”, “Plan”)
- Center: Glass Hub
- Right: Secondary view (e.g., “History”, “Grocery List”)

---

## 8. Theme Switching Per Section

Theme changes are subtle and accent-driven.

### Global Rules
- One neutral base palette (grays, whites)
- One accent color per section
- Minimal background gradient shifts

### Accent Usage
- Hub glow
- Connector tint
- Key buttons
- Charts and highlights
- Header chips

### Example Accents
- Calories: Warm coral / amber
- Meal Plan: Fresh green
- Recipes: Soft orange
- Insights: Cool blue / purple

### Transition Behavior
- Crossfade between sections
- Slight parallax on background
- No full recoloring of UI

---

## 9. UI Details That Sell “Apple Liquid Glass”

- Large corner radii throughout
- Layered depth with floating cards
- Confident, large typography
- Minimal font weight changes
- Simple outline icons with consistent stroke
- Subtle haptic feedback on:
  - Hub open/close
  - Bubble focus
  - Selection confirm
  - Double-press Home

---

## 10. Radial Menu Variants

### Variant A: Arc Fan (Recommended)
- 120–160° arc above hub
- Best for one-handed use

---

## Summary

The Glass Hub navigation combines:
- Radial interaction
- Liquid glass materials
- Contextual navigation modes
- Subtle theme switching
- Apple-style motion and restraint

It provides a playful yet native-feeling way to navigate between major sections without overwhelming the UI.

---

