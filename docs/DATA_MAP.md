# Data Map

## 1. Data Entities / Tables
- **users** – Stores user account information and current credit balance.
  - `id` (BIGINT, primary key, auto-increment)
  - `wechat_openid` (VARCHAR(255), unique, not null) - For linking WeChat account
  - `nickname` (VARCHAR(255))
  - `avatar_url` (VARCHAR(1024))
  - `credits` (INT, not null, default 0) - User's "胶卷" balance
  - `created_at` (TIMESTAMP, default CURRENT_TIMESTAMP)
  - `updated_at` (TIMESTAMP, default CURRENT_TIMESTAMP on update)

- **templates** – Stores available AI style templates.
  - `id` (INT, primary key, auto-increment)
  - `name` (VARCHAR(255), not null)
  - `description` (TEXT)
  - `preview_image_url` (VARCHAR(1024), not null) - URL to the example image
  - `credit_cost` (INT, not null) - How many credits this template costs to use
  - `is_active` (BOOLEAN, not null, default true) - To enable/disable templates
  - `created_at` (TIMESTAMP, default CURRENT_TIMESTAMP)

- **transactions** – Logs all credit changes for auditing and user history.
  - `id` (BIGINT, primary key, auto-increment)
  - `user_id` (BIGINT, not null, foreign key to users.id)
  - `type` (ENUM('purchase', 'generation'), not null)
  - `amount` (INT, not null) - Positive for purchases, negative for generation
  - `description` (VARCHAR(255)) - e.g., "Purchased 200 credits" or "Used 'Dusty Rose' template"
  - `external_payment_id` (VARCHAR(255), nullable) - ID from WeChat Pay/Apple
  - `related_template_id` (INT, nullable, foreign key to templates.id)
  - `created_at` (TIMESTAMP, default CURRENT_TIMESTAMP)

## 2. Entity Relationships
- A **user** can have many **transactions** (1-to-N).
- A **transaction** of type `generation` is related to one **template** (N-to-1).

## 3. Data Flow / Lifecycle
1.  **User Creation:** A new `users` record is created upon first WeChat login.
2.  **Purchase Flow:** User payment triggers the creation of a `transactions` record with `type='purchase'`, and the `credits` field in the `users` table is incremented.
3.  **Generation Flow:** User generating an image triggers the creation of a `transactions` record with `type='generation'`, and the `credits` field in the `users` table is decremented. The `related_template_id` is logged.

## 4. Data Constraints and Validation
- `users.wechat_openid` must be unique.
- `users.credits` cannot be negative. The application layer must check for sufficient funds before creating a `generation` transaction.
- Foreign key constraints must be enforced between `transactions` and `users`/`templates`.

## 5. Sample Data (JSON Representation)
```json
{
  "user": {
    "id": 101,
    "wechat_openid": "oGZUI0abcdefghijklmnopq",
    "nickname": "月亮",
    "avatar_url": "https://.../avatar.jpg",
    "credits": 150,
    "created_at": "2024-05-20T10:00:00Z"
  },
  "template": {
    "id": 5,
    "name": "Fairy Tale Dream",
    "preview_image_url": "https://cdn.45ai.com/templates/fairy_tale.jpg",
    "credit_cost": 15,
    "is_active": true
  },
  "transaction_purchase": {
    "id": 2001,
    "user_id": 101,
    "type": "purchase",
    "amount": 200,
    "description": "Purchased 200 胶卷",
    "external_payment_id": "wx20240520...",
    "created_at": "2024-05-20T11:00:00Z"
  },
  "transaction_generation": {
    "id": 2002,
    "user_id": 101,
    "type": "generation",
    "amount": -15,
    "description": "Used 'Fairy Tale Dream' template",
    "related_template_id": 5,
    "created_at": "2024-05-20T11:05:00Z"
  }
}