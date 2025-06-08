# Test Plan

## 1. Testing Strategy
- **Unit Tests:** Go's built-in testing package for backend logic. Vitest for frontend utility functions and state management logic.
- **Integration Tests:** API endpoints will be tested using a Go testing suite to verify interactions between the Go API and the MySQL database.
- **End-to-End (E2E) Tests:** Using WeChat's Miniprogram automation tools to simulate key user flows.
- **Visual Regression Testing:** A tool like Percy or Applitools could be considered to automatically catch unintended UI changes, ensuring the refined visual style remains consistent.
- **Manual QA:** Critical for verifying the feel of animations, transitions, and overall UX on physical iOS devices and within WeChat.

## 2. Test Environment
- A dedicated `staging` environment on WeChat Cloud Hosting, connected to a staging database.
- The staging environment will use a sandboxed ComfyUI instance and test modes for content moderation and payment APIs.
- Environment variables will distinguish between `development`, `staging`, and `production`.

## 3. Test Cases by Feature

**Feature: Image Generation & UX**
- **Test GEN-01:** Successful Flow - Test the entire generation flow with a valid selfie and sufficient credits. - **Expected:** Credits are deducted, the correct soothing loading animation appears, and images are revealed with a staggered fade-in.
- **Test GEN-02:** Invalid Photo (Face Detect) - Upload an image of a landscape. - **Expected:** A softly styled error alert appears with the message "Please upload a photo with a face."
- **Test GEN-03:** Invalid Photo (Content Safety) - Upload a flagged image. - **Expected:** A generic error "Upload failed. Please try a different photo." is displayed.
- **Test GEN-04:** Insufficient Credits - Attempt generation when credits are too low. - **Expected:** A modal prompts the user to purchase more credits.

**Feature: User & Credit System**
- **Test USR-01:** First-time Login - A new user authorizes via WeChat. - **Expected:** A new user record is created, and the user is transitioned to the home screen with the correct fade/slide animation.
- **Test USR-02:** Purchase Credits - Successfully complete a purchase. - **Expected:** `transactions` table is updated, and the credit balance on the Profile page updates with a number-scroll animation.
- **Test USR-03:** View Transaction History - Navigate to the history page. - **Expected:** A list of all credit purchases and expenses is displayed with clean typography and spacing.

**Feature: UI & Motion**
- **Test UI-01:** Page Transitions - Navigate between the Home, Profile, and a Template Detail screen. - **Expected:** All page transitions use the specified "Subtle Fade + Vertical Slide" animation.
- **Test UI-02:** Card & Button Press State - Tap and hold on a template card and a primary button. - **Expected:** The elements scale and change shadow according to the Style Guide's micro-interaction specs.
- **Test UI-03:** Icon States - Tap on the main navigation icons (e.g., Home, Profile). - **Expected:** The inactive icon is a thin outline; the active icon fills with the primary color.

## 4. Non-Functional Tests
- **Performance:** Load test key API endpoints. Measure frame rate (FPS) during animations on target devices to ensure they are smooth (>55 FPS).
- **Security:** Run vulnerability scans for common web security issues.

## 5. Test Data
- A pool of test user accounts with varying credit balances.
- A set of test images: valid selfies, group photos, landscapes, and images known to trigger content moderation filters.
- A dummy set of template data in the test database.