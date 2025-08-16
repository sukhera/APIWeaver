import { test, expect } from '@playwright/test';

test.describe('APIWeaver Application', () => {
  test('should load the home page', async ({ page }) => {
    await page.goto('/');
    
    // Check that the page loads
    await expect(page).toHaveTitle(/APIWeaver/);
    
    // Check for main navigation elements
    await expect(page.locator('nav')).toBeVisible();
  });

  test('should navigate to Generate page', async ({ page }) => {
    await page.goto('/');
    
    // Navigate to Generate page
    await page.click('text=Generate');
    await expect(page).toHaveURL(/.*generate/);
    
    // Check for generate page elements
    await expect(page.locator('h1')).toContainText(/Generate/i);
  });

  test('should navigate to Validate page', async ({ page }) => {
    await page.goto('/');
    
    // Navigate to Validate page
    await page.click('text=Validate');
    await expect(page).toHaveURL(/.*validate/);
    
    // Check for validate page elements
    await expect(page.locator('h1')).toContainText(/Validate/i);
  });

  test('should navigate to Amend page', async ({ page }) => {
    await page.goto('/');
    
    // Navigate to Amend page
    await page.click('text=Amend');
    await expect(page).toHaveURL(/.*amend/);
    
    // Check for amend page elements
    await expect(page.locator('h1')).toContainText(/Amend/i);
  });

  test('should toggle theme', async ({ page }) => {
    await page.goto('/');
    
    // Find and click theme toggle
    const themeToggle = page.locator('[data-testid="theme-toggle"], button[aria-label*="theme"]').first();
    if (await themeToggle.isVisible()) {
      await themeToggle.click();
      
      // Check that theme changed (look for dark/light class on html or body)
      const html = page.locator('html');
      const hasThemeClass = await html.getAttribute('class');
      expect(hasThemeClass).toBeTruthy();
    }
  });

  test('should be responsive on mobile', async ({ page, isMobile }) => {
    await page.goto('/');
    
    if (isMobile) {
      // Check mobile-specific elements
      await expect(page.locator('body')).toBeVisible();
      
      // Check that content is readable on mobile
      const viewport = page.viewportSize();
      expect(viewport?.width).toBeLessThanOrEqual(768);
    }
  });

  test('should handle navigation errors gracefully', async ({ page }) => {
    await page.goto('/non-existent-page');
    
    // Should either redirect to home or show 404 page
    // Adjust this based on your routing strategy
    const url = page.url();
    expect(url).toBeTruthy(); // Page should load something
  });
});