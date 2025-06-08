# User Experience Flows

## Overview
The app follows a minimalist, visually-driven design philosophy focused on intuitive, single-purpose interactions. The experience is heavily inspired by the elegant simplicity of 醒图 (Xingtu), prioritizing aesthetic appeal and ease of use. The UI is clean and airy, with soft colors and fluid animations to create a serene and premium user journey.

## Flow 1: Core Image Generation
**Goal:** User wants to select a style and generate an AI photo of themselves.
**Actors:** User (Logged in)
**Steps:**
1.  **Home:** User opens the app and lands on the Home screen, displaying the Template Gallery with a gentle fade-in animation.
2.  **Browse & Select:** User scrolls through template cards. Tapping a card causes it to lift subtly, providing tactile feedback before navigating to a detail view or modal.
3.  **Initiate Generation:** User taps the primary "Generate with this Template" button. The button provides a press-in animation.
4.  **Upload:** The system prompts the user to upload a photo, using the native OS file picker.
5.  **Validation:** The system validates the image in the background.
    - **On Success:** Proceeds seamlessly.
    - **On Failure:** An alert with soft styling appears (e.g., "Please upload a clear photo of a single person") and the user can try again.
6.  **Credit Deduction:** The system deducts the required credits from the user's account.
7.  **Generate & Wait:** The UI transitions to a waiting screen. A soothing, branded animation (e.g., a pulsating, soft-colored blob) is displayed with text like "AI is creating your look... ✨".
8.  **Display Results:** After generation, the results screen appears, with the 1-4 generated images fading in sequentially with a slight upward drift.
9.  **Save:** User can tap on any result to view it larger, and then tap a "Save to Album" button. A confirmation toast appears ("Saved successfully!").
10. **Return:** User can go back to the Home screen via a soft-edged back icon.

**UX Notes:** Motion is key. Transitions should be fluid and gentle. The waiting animation is crucial for managing perceived wait time and reinforcing the brand's soft aesthetic.

## Flow 2: User Onboarding & Login
**Goal:** A new user registers or an existing user logs in.
**Actors:** New or Returning User
**Steps:**
1.  **Launch:** User opens the app. The launch screen briefly shows the app logo, then fades into the login prompt.
2.  **Login Prompt:** The app prompts for WeChat login. The screen is minimal, featuring the value proposition and a primary "Login with WeChat" button.
3.  **Authorize:** User taps the button and authorizes the app in the native WeChat prompt.
4.  **Complete:** The system creates/logs in the user. The user is transitioned to the Home screen via a gentle fade/slide-up animation.

**UX Notes:** The onboarding must be frictionless. The goal is to move the user into the core experience in the most elegant way possible.

## Flow 3: Credit Recharge
**Goal:** User wants to purchase more credits ("胶卷").
**Actors:** User (Logged in)
**Steps:**
1.  **Initiate Purchase:** User navigates to the "Profile" tab (via a thin-line icon) and taps on their credit balance, or is prompted after attempting generation with insufficient funds.
2.  **Select Pack:** The user is shown a list of credit packs presented in clean cards with clear pricing.
3.  **Confirm:** User selects a pack. The selected card might subtly expand or show a checkmark. They tap a "Purchase" button.
4.  **Payment:** The app initiates the standard WeChat Pay or Apple In-App Purchase flow.
5.  **Complete:** Upon successful payment, the UI updates the credit balance with a subtle number-scrolling animation. A confirmation message is shown.
6.  **Return:** The user is returned to the previous screen.

**UX Notes:** The purchase screen should feel secure and trustworthy, using familiar UI patterns and clear typography.