# Style Guide (v2.0 - Refined Aesthetic)

## 1. Design Philosophy
The UI should feel like a premium, serene digital space. We are aiming for a **minimalist, light, and soft feminine aesthetic**. Key principles are ample whitespace, subtle gradients, soft shadows, and fluid, non-jarring animations. The experience should be calming and visually pleasing, avoiding loud colors and overly technical elements. Think "digital skincare" more than "tech tool."

## 2. Brand Palette
A muted, warm, and sophisticated palette inspired by Morandi tones and soft pastels.

| Token | HEX | Usage |
| ----- | --- | ----- |
| --color-primary | #E89B93 | **Dusty Rose.** Main interactive elements: buttons, links, active icons. |
| --color-secondary | #F3D9D7 | **Powder Pink.** Accents, tags, subtle backgrounds, loading glows. |
| --color-bg | #FCFBF9 | **Alabaster.** App background, feels softer than pure white. |
| --color-surface | #FFFFFF | **White.** Card and input backgrounds, for a clean contrast on the alabaster. |
| --color-text | #4A4A4A | **Charcoal.** Main body text, for readability without harshness. |
| --color-text-subtle | #9B9B9B | **Stone Grey.** Captions, placeholder text, disabled states. |
| --color-gradient | `linear-gradient(to top right, #F3D9D7, #E89B93)` | Optional use on banners or special UI elements for a premium feel. |

## 3. Typography
The focus is on clarity and elegance. Generous line-height is crucial for the airy feel.

| Level | Font | Size / Line-height | Weight |
| ----- | ---- | ------------------ | ------ |
| H1 | PingFang SC / Inter | 26 / 38 px | 600 (Semibold) |
| H2 | PingFang SC / Inter | 20 / 30 px | 600 (Semibold) |
| Body | PingFang SC / Inter | 15 / 26 px | 400 (Regular) |
| Caption | PingFang SC / Inter | 12 / 20 px | 400 (Regular) |

## 4. Spacing & Sizing
Generous spacing is the foundation of our minimalist aesthetic.

- **Grid unit:** 4 px
- **Key Spacings:** Use 8px, 12px, 16px, 24px, 32px. Avoid random values.
- **Page Padding:** 20px horizontal padding on all main screens.
- **Card Padding:** 16px to 24px internally, depending on content.
- **Layout Principle:** Prioritize whitespace. Do not crowd elements together.

## 5. Components(we use wotui in this project as part of UniBest)
### 5.1 Button
- **Primary:** Background: `--color-primary`. Text: `#FFF`. Radius: `24px` (pill-shaped). Shadow: `0 4px 12px rgba(232, 155, 147, 0.3)`. A soft glow from its own color.
- **Secondary:** Border: `1px solid --color-primary`. Background: `transparent`. Text: `--color-primary`. Radius: `24px`. No shadow.

### 5.2 Card (e.g., Template Card)
- **Background:** `--color-surface` (#FFFFFF)
- **Border:** `1px solid #F0F0F0` (subtle definition against the alabaster background)
- **Border-radius:** 16 px
- **Shadow:** A very soft, diffuse shadow: `0 4px 24px rgba(74, 74, 74, 0.08)`

### 5.3 Iconography
- **Style:** Thin-line, minimalist, with a consistent stroke weight (e.g., 1.5px).
- **Details:** Softly rounded corners and caps. Avoid sharp angles.
- **Active State:** The active icon should fill with `--color-primary`.

## 6. Motion & Animation
Motion is a core part of the "design sense." It must be fluid, gentle, and meaningful.

- **Easing:** A custom cubic-bezier `cubic-bezier(0.6, 0.05, 0.4, 1)` for most transitions, creating a graceful 'ease in and out' effect.
- **Page Transitions:** **Subtle Fade + Vertical Slide.** New pages gently fade in while sliding up slightly (`transform: translateY(16px)` to `translateY(0)`). Duration: 400ms.
- **Loading Animation:** Instead of a spinner, use a **softly pulsing blob or series of circles**, animated with the `--color-secondary` and `--color-primary` tones. The animation should be organic and soothing.
- **Image Reveal:** Generated images should not just appear. Use a **staggered fade-in** with a slight upward drift. Each image appears 100ms after the previous one. `opacity: 0 -> 1`, `transform: translateY(10px) -> 0`.
- **Micro-interactions:**
    - **Button Tap:** On press, the button scales down slightly (`scale(0.97)`) and its shadow becomes smaller.
    - **Card Tap:** On press, the card lifts up and forwards (`scale(1.02)`, `transform: translateY(-4px)`) with a slightly stronger shadow to indicate it's active.

## 7. Accessibility
- **Color-contrast:** Ensure `--color-text` on `--color-surface` and other key combinations meet a minimum 4.5:1 ratio.
- **Tap Targets:** Minimum 48x48px tappable area for all interactive elements to ensure ease of use.
- **Focus States:** A subtle, 2px outline using `--color-secondary` for any keyboard-accessible elements in the future.