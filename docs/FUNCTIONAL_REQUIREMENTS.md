# Functional Requirements Specification

## 1. User Account & Authentication Module

| ID          | Requirement                   | Description                                                                                                                                                |
| :---------- | :---------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **FR-U-01** | **User Registration & Login** | The system must allow new users to register via email or social login and allow existing users to log in securely.                                         |
| **FR-U-02** | **Password Recovery**         | The system must provide a mechanism for users to securely reset their password via a registered email address.                                             |
| **FR-U-03** | **Profile Management**        | The system must allow logged-in users to update their profile information (name, shipping address, billing details) and view their order history.          |
| **FR-U-04** | **Role-Based Access Control** | The system must enforce access rules based on defined roles: **Customer**, **Administrator**, **Moderator**, **Content Manager**, and **Order Processor**. |

---

## 2. Product Catalog & Management Module

| ID          | Requirement                     | Description                                                                                                                                                       |
| :---------- | :------------------------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **FR-P-01** | **Product Creation**            | The system must allow Content Managers to create a new product entry, including name, description, images, and category assignment.                               |
| **FR-P-02** | **Hierarchical Categorization** | The system must support a hierarchical product category structure (e.g., Electronics > Phones > Android) for effective browsing.                                  |
| **FR-P-03** | **Product Variants**            | The system must allow products to have multiple variants (e.g., size, color) with separate stock keeping units (SKUs) and potentially separate pricing/inventory. |
| **FR-P-04** | **Inventory Tracking**          | The system must accurately track the quantity of each product SKU and prevent an order from being placed if inventory is zero.                                    |
| **FR-P-05** | **Multi-Currency Pricing**      | The system must allow prices to be defined in a base currency and display converted prices to users based on their region.                                        |

---

## 3. Search & Discovery Module

| ID          | Requirement                    | Description                                                                                                                              |
| :---------- | :----------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------- |
| **FR-S-01** | **Full-Text Search**           | The system must allow users to search the entire product catalog based on keywords in the product name, description, and attributes.     |
| **FR-S-02** | **Faceted Search / Filtering** | Search results must be filterable by product attributes (e.g., brand, price range, color, size) and category.                            |
| **FR-S-03** | **Sorting**                    | Users must be able to sort search results and category pages by relevance, price (low to high/high to low), and average customer rating. |
| **FR-S-04** | **Search Auto-Completion**     | The search bar must provide real-time suggestions based on popular queries and existing product names as the user types.                 |

---

## 4. Product Review & Rating Module

| ID          | Requirement                     | Description                                                                                                                                                      |
| :---------- | :------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **FR-R-01** | **Review Submission**           | Only authenticated users who have **purchased the product** must be allowed to submit a written review and a star rating (1-5).                                  |
| **FR-R-02** | **Rating Aggregation**          | The system must asynchronously calculate and display the average star rating for each product, updating whenever a new review is submitted.                      |
| **FR-R-03** | **Review Moderation**           | A Moderator role must be able to review, approve, or reject submitted reviews before they are visible to the public.                                             |
| **FR-R-04** | **Review Voting (Helpfulness)** | The system must allow all logged-in users to mark a review as "Helpful" or "Unhelpful," with logic in place to prevent a single user from voting multiple times. |
