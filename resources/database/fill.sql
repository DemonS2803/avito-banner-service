

INSERT INTO features (id, name)
VALUES
    (1, 'Product Listing'),
    (2, 'Product Details Page'),
    (3, 'Add to Cart Functionality'),
    (4, 'Checkout Process'),
    (5, 'User Account Management'),
    (6, 'Order Tracking'),
    (7, 'Product Reviews and Ratings'),
    (8, 'Payment Gateway Integration'),
    (9, 'Shipping Options'),
    (10, 'Promotional Discounts and Coupons');

INSERT INTO tags (id, name)
VALUES
    (1, 'Fashion'),
    (2, 'Technology'),
    (3, 'Food'),
    (4, 'Travel');

INSERT INTO user_banners (id, content, is_active, feature_id)
VALUES
    (1, '{
  "title": "New Arrivals in Fashion",
  "text": "Discover the latest trends and styles for the season.",
  "url": "https://example.com/new-fashion-arrivals"
}', true, 1),
    (2, '{
  "title": "Tech Deals of the Week",
  "text": "Explore amazing discounts on the hottest tech gadgets.",
  "url": "https://example.com/tech-deals"
}', true, 2),
    (3, '{
  "title": "Culinary Masterclass",
  "text": "Learn cooking tips and tricks from top chefs.",
  "url": "https://example.com/culinary-masterclass"
}', true, 3),
    (4, '{
  "title": "Plan Your Dream Vacation",
  "text": "Find exclusive travel deals to your dream destinations.",
  "url": "https://example.com/dream-vacation-deals"
}', true, 4),
    (5, '{
  "title": "Fitness Challenge",
  "text": "Join our fitness challenge and get in shape!",
  "url": "https://example.com/fitness-challenge"
}', true, 5),
    (6, '{
  "title": "Home Decor Inspiration",
  "text": "Get inspired with the latest home decor trends.",
  "url": "https://example.com/home-decor-inspiration"
}', true, 6),
    (7, '{
  "title": "Gaming News Update",
  "text": "Stay updated with the latest gaming news and releases.",
  "url": "https://example.com/gaming-news"
}', false, 7),
    (8, '{
  "title": "Wellness Retreat",
  "text": "Recharge your mind and body with our wellness retreat.",
  "url": "https://example.com/wellness-retreat"
}', true, 8),
    (9, '{
  "title": "Bookworm''s Paradise",
  "text": "Discover your next favorite book with our recommendations.",
  "url": "https://example.com/bookworm-paradise"
}', true, 9),
    (10, '{
  "title": "DIY Ideas for Home",
  "text": "Get creative with our DIY home improvement projects.",
  "url": "https://example.com/diy-home-ideas"
}', true, 10);




INSERT INTO user_banners_tags (banner_id, tag_id) VALUES
                                                      (1, 1),
                                                      (1, 2),
                                                      (2, 2),
                                                      (2, 3),
                                                      (3, 3),
                                                      (3, 4),
                                                      (4, 4),
                                                      (4, 2),
                                                      (5, 3),
                                                      (5, 1),
                                                      (6, 1),
                                                      (6, 3),
                                                      (7, 2),
                                                      (7, 4),
                                                      (8, 1),
                                                      (8, 4),
                                                      (9, 1),
                                                      (9, 1),
                                                      (10, 3),
                                                      (10, 2);
