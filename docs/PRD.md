# Product Requirements Document: 45AI

## 1. Executive Summary
45AI is a cross-platform AI image generation app for iOS and WeChat, designed to solve the problem of overly complex and low-quality AI art tools. It targets young, aesthetically-minded women in China who are accustomed to polished photo editing apps. The core solution is a one-click experience: users upload a single selfie and select a professionally designed template to receive beautiful, high-quality AI-generated portraits instantly, monetized through a simple credit-based system.

## 2. Target Audience & Problem Statement
- **Primary User:** Young women in China (18-35) who are active on social media, use apps like 醒图 (Xingtu) and 美图秀秀 (Meitu), and value aesthetic quality and convenience.
- **Problem:** Existing AI image generators are often cumbersome, requiring multiple photo uploads, complex text prompts, or long waiting times. The resulting styles frequently fail to meet the user's aesthetic expectations for shareable, high-quality portraits.

## 3. Core Features & User Stories

### 3.1 Feature: Template-Based Image Generation
- **As a** user, **I want to** browse a gallery of pre-made visual styles (templates), **so that** I can easily find an aesthetic I love without needing to describe it.
- **As a** user, **I want to** upload a single selfie and select a template, **so that** the AI can generate a stylized image for me with one click.
- **As a** user, **I want to** save the generated high-quality images to my phone's photo album, **so that** I can easily share them on my social media accounts.

### 3.2 Feature: Credit ("胶卷") Monetization System
- **As a** user, **I want to** clearly see how many credits each template costs, **so that** I can make an informed decision before generating an image.
- **As a** user, **I want to** purchase packs of credits through a simple and secure payment process (WeChat Pay/Apple IAP), **so that** I can continue using the service when I run out.
- **As a** user, **I want to** view my purchase and spending history, **so that** I can track my consumption on the platform.

### 3.3 Feature: Content Safety & Moderation
- **As a** user, **I want to** trust that the platform is a safe environment, **so that** I can upload my photo without being exposed to inappropriate content from others.
- **As a** platform operator, **I want** all user-uploaded images to be automatically scanned for unsafe content, **so that** the platform remains compliant with regulations and protects its brand reputation.

## 4. Non-Functional Requirements (NFRs)
- **Performance:**
  - Backend API endpoints (non-generation) must respond in under 200ms.
  - The P95 end-to-end image generation time (from upload confirmation to image display) must be under 30 seconds.
  - Frontend animations must maintain a smooth frame rate (>55 FPS) on target devices.
- **Scalability:**
  - The system must be architected to support 1,000 concurrent image generation requests in v1.
- **Security:**
  - All user data (especially PII) must be encrypted at rest in the database.
  - All network communication must use HTTPS.
  - All user-uploaded images must be processed by a content safety moderation API before being sent to the generation model.
- **Accessibility:**
  - Frontend must provide minimum tap targets of 48x48px for all interactive elements.
  - Text color contrast must meet WCAG 2.1 AA standards for readability.

## 5. Scope: In & Out
- **In Scope for v1:**
  - WeChat Mini Program and native iOS application.
  - User authentication via WeChat Login.
  - A credit-based ("胶卷") monetization system.
  - Image generation based on a curated library of pre-defined templates.
- **Out of Scope for v1:**
  - User-defined text prompts.
  - Social features (e.g., in-app profiles, sharing, following).
  - Storing user-generated images or upload history long-term.
  - A web or Android version.

## 6. Assumptions & Dependencies
- **Assumptions:**
  - Users have a strong preference for simplicity and are willing to pay for a high-quality, convenient experience.
  - The aesthetic quality and trendiness of the templates are the primary drivers for user adoption and retention.
- **Dependencies:**
  - Relies on a self-hosted ComfyUI API deployed on GCP for the core image generation task.
  - Relies on a third-party Content Safety API (e.g., Tencent Cloud, Volcano Engine) for image moderation.
  - Relies on platform-native services (WeChat SDK, Apple StoreKit) for authentication and In-App Purchases.