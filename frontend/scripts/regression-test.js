/**
 * 45AI Cross-Platform Regression Test Suite
 * Tests all UI/UX refinements across H5, WeChat Mini-Program, and iOS
 */

const { execSync } = require('child_process')
const { existsSync, writeFileSync, readFileSync } = require('fs')
const { join } = require('path')

console.log('🧪 Starting 45AI Cross-Platform Regression Test Suite...')

const testResults = {
  timestamp: new Date().toISOString(),
  platform: process.env.UNI_PLATFORM || 'h5',
  totalTests: 0,
  passedTests: 0,
  failedTests: 0,
  results: [],
}

// Test Categories
const testCategories = {
  UI_REFINEMENTS: 'UI/UX Refinements',
  NAVIGATION: 'Navigation & Layout',
  AUTHENTICATION: 'Authentication Flow',
  GENERATION: 'Image Generation',
  CREDITS: 'Credit System',
  VISUAL_STYLE: 'Visual Style Consistency',
}

/**
 * Test Runner Function
 */
async function runTest(testName, testFn, category = testCategories.UI_REFINEMENTS) {
  testResults.totalTests++

  try {
    console.log(`\n🔍 Running: ${testName}`)
    await testFn()

    testResults.passedTests++
    testResults.results.push({
      name: testName,
      category,
      status: 'PASSED',
      timestamp: new Date().toISOString(),
    })

    console.log(`✅ PASSED: ${testName}`)
  } catch (error) {
    testResults.failedTests++
    testResults.results.push({
      name: testName,
      category,
      status: 'FAILED',
      error: error.message,
      timestamp: new Date().toISOString(),
    })

    console.log(`❌ FAILED: ${testName}`)
    console.log(`   Error: ${error.message}`)
  }
}

/**
 * Visual Style Tests
 */
async function testSofterVisualStyle() {
  // Check if new color variables are properly applied
  const variablesFile = readFileSync(join(process.cwd(), 'src/style/variables.scss'), 'utf8')

  if (!variablesFile.includes('$color-primary: #F4A5A8')) {
    throw new Error('Softer primary color not applied')
  }

  if (!variablesFile.includes('$color-accent-peach: #F7E6D7')) {
    throw new Error('Soft accent colors not defined')
  }

  if (!variablesFile.includes('$shadow-subtle: 0 1px 4px')) {
    throw new Error('Subtle shadows not implemented')
  }
}

async function testBottomNavigationMinimalist() {
  const pagesConfig = readFileSync(join(process.cwd(), 'pages.config.ts'), 'utf8')

  // Check for minimalist design (text-only, no icons)
  if (pagesConfig.includes('iconPath')) {
    throw new Error('Bottom navigation still using icons (should be text-only)')
  }

  if (!pagesConfig.includes("height: '56px'")) {
    throw new Error('Bottom navigation height not updated to 56px')
  }
}

async function testTemplateCardConsistency() {
  const templateCard = readFileSync(
    join(process.cwd(), 'src/components/TemplateCard/TemplateCard.vue'),
    'utf8',
  )

  // Check for 3:4 aspect ratio
  if (!templateCard.includes('aspect-ratio: 3/4')) {
    throw new Error('Template cards not using consistent 3:4 aspect ratio')
  }

  // Check for consistent navigation
  if (!templateCard.includes('navigateTo') || !templateCard.includes('template-detail')) {
    throw new Error('Template card navigation not properly implemented')
  }
}

/**
 * Authentication Flow Tests
 */
async function testTokenRefreshFlow() {
  const requestUtils = readFileSync(join(process.cwd(), 'src/utils/request.ts'), 'utf8')

  // Check for improved token refresh implementation
  if (!requestUtils.includes('useUserStore')) {
    throw new Error('Token refresh not using real user store')
  }

  if (!requestUtils.includes('refreshUserToken')) {
    throw new Error('Real token refresh API not integrated')
  }
}

/**
 * Navigation Tests
 */
async function testPageExistenceValidation() {
  // Check that navigation targets exist
  const galleryPage = existsSync(join(process.cwd(), 'src/pages/gallery/index.vue'))
  const profilePage = existsSync(join(process.cwd(), 'src/pages/profile/index.vue'))
  const purchasePage = existsSync(join(process.cwd(), 'src/pages/purchase/index.vue'))
  const templateDetailPage = existsSync(join(process.cwd(), 'src/pages/template-detail/index.vue'))

  if (!galleryPage) throw new Error('Gallery page missing')
  if (!profilePage) throw new Error('Profile page missing')
  if (!purchasePage) throw new Error('Purchase page missing')
  if (!templateDetailPage) throw new Error('Template detail page missing')
}

/**
 * Data Consistency Tests
 */
async function testTemplateDataConsistency() {
  const galleryPage = readFileSync(join(process.cwd(), 'src/pages/gallery/index.vue'), 'utf8')
  const templateDetailPage = readFileSync(
    join(process.cwd(), 'src/pages/template-detail/index.vue'),
    'utf8',
  )

  // Both pages should use template store
  if (!galleryPage.includes('useTemplateStore')) {
    throw new Error('Gallery page not using template store')
  }

  if (!templateDetailPage.includes('useTemplateStore')) {
    throw new Error('Template detail page not using template store')
  }
}

/**
 * Purchase Flow Tests
 */
async function testPurchasePageGrid() {
  const purchasePage = readFileSync(join(process.cwd(), 'src/pages/purchase/index.vue'), 'utf8')

  // Check for 2x3 grid layout
  if (!purchasePage.includes('grid-template-columns: repeat(2, 1fr)')) {
    throw new Error('Purchase page not using 2x3 grid layout')
  }

  // Check for proper radio button behavior in selectPack method
  if (!purchasePage.includes('this.selectedPack = pack')) {
    throw new Error('Purchase page not implementing radio selection behavior')
  }

  // Verify deselection logic is properly contained (only after purchase)
  const selectPackSection = purchasePage.substring(
    purchasePage.indexOf('selectPack(pack)'),
    purchasePage.indexOf('async handlePurchase()'),
  )

  if (selectPackSection.includes('selectedPack = null')) {
    throw new Error(
      'Purchase page allows deselection in selectPack method (should be radio behavior)',
    )
  }
}

/**
 * Profile Page Tests
 */
async function testProfileUsernameOverflow() {
  const profilePage = readFileSync(join(process.cwd(), 'src/pages/profile/index.vue'), 'utf8')

  // Check for improved overflow handling
  if (!profilePage.includes('showFullUsername')) {
    throw new Error('Profile page username overflow not implemented')
  }

  if (profilePage.includes('max-width: 200px')) {
    throw new Error('Profile page still using fixed max-width (should be responsive)')
  }
}

/**
 * Platform-Specific Tests
 */
async function testH5Compatibility() {
  if (process.env.UNI_PLATFORM === 'h5') {
    const generationStore = readFileSync(join(process.cwd(), 'src/store/generation.ts'), 'utf8')

    // Check for H5-specific upload implementation
    if (!generationStore.includes('XMLHttpRequest') || !generationStore.includes('FormData')) {
      throw new Error('H5-specific file upload not implemented')
    }
  }
}

/**
 * Component Style Tests
 */
async function testComponentSoftening() {
  const componentsScss = readFileSync(join(process.cwd(), 'src/style/components.scss'), 'utf8')

  // Check for softer button styles
  if (!componentsScss.includes('linear-gradient(135deg, $color-primary')) {
    throw new Error('Button gradients not softened')
  }

  // Check for softer shadows
  if (!componentsScss.includes('$shadow-button')) {
    throw new Error('Softer button shadows not applied')
  }

  // Check for soft enhancement classes
  if (!componentsScss.includes('.soft-surface')) {
    throw new Error('Soft enhancement classes not added')
  }
}

/**
 * Run All Tests
 */
async function runAllTests() {
  console.log('🎯 Testing UI/UX Refinements...')

  // Visual Style Tests
  await runTest('Softer Visual Style Applied', testSofterVisualStyle, testCategories.VISUAL_STYLE)
  await runTest(
    'Bottom Navigation Minimalist Design',
    testBottomNavigationMinimalist,
    testCategories.UI_REFINEMENTS,
  )
  await runTest(
    'Template Card Consistency',
    testTemplateCardConsistency,
    testCategories.UI_REFINEMENTS,
  )
  await runTest('Component Style Softening', testComponentSoftening, testCategories.VISUAL_STYLE)

  // Authentication Tests
  await runTest('Token Refresh Flow Fixed', testTokenRefreshFlow, testCategories.AUTHENTICATION)

  // Navigation Tests
  await runTest('Page Existence Validation', testPageExistenceValidation, testCategories.NAVIGATION)

  // Data Consistency Tests
  await runTest(
    'Template Data Consistency',
    testTemplateDataConsistency,
    testCategories.UI_REFINEMENTS,
  )

  // Purchase Flow Tests
  await runTest('Purchase Page Grid Layout', testPurchasePageGrid, testCategories.UI_REFINEMENTS)

  // Profile Tests
  await runTest(
    'Profile Username Overflow Fix',
    testProfileUsernameOverflow,
    testCategories.UI_REFINEMENTS,
  )

  // Platform-Specific Tests
  await runTest('H5 Platform Compatibility', testH5Compatibility, testCategories.UI_REFINEMENTS)

  // Generate Report
  generateReport()
}

/**
 * Generate Test Report
 */
function generateReport() {
  const reportPath = join(process.cwd(), 'test-reports', `regression-report-${Date.now()}.json`)

  // Ensure reports directory exists
  execSync('mkdir -p test-reports', { stdio: 'inherit' })

  // Calculate success rate
  const successRate =
    testResults.totalTests > 0
      ? ((testResults.passedTests / testResults.totalTests) * 100).toFixed(1)
      : 0

  const report = {
    ...testResults,
    successRate: `${successRate}%`,
    summary: {
      total: testResults.totalTests,
      passed: testResults.passedTests,
      failed: testResults.failedTests,
      platform: testResults.platform,
    },
  }

  writeFileSync(reportPath, JSON.stringify(report, null, 2))

  console.log('\n📊 TEST REPORT SUMMARY')
  console.log('=====================')
  console.log(`Platform: ${testResults.platform}`)
  console.log(`Total Tests: ${testResults.totalTests}`)
  console.log(`Passed: ${testResults.passedTests}`)
  console.log(`Failed: ${testResults.failedTests}`)
  console.log(`Success Rate: ${successRate}%`)
  console.log(`Report saved to: ${reportPath}`)

  // List failed tests
  if (testResults.failedTests > 0) {
    console.log('\n❌ Failed Tests:')
    testResults.results
      .filter((r) => r.status === 'FAILED')
      .forEach((r) => console.log(`   - ${r.name}: ${r.error}`))
  }

  console.log('\n🎉 Regression test suite completed!')

  // Exit with error code if tests failed
  if (testResults.failedTests > 0) {
    process.exit(1)
  }
}

// Run the test suite
runAllTests().catch((error) => {
  console.error('💥 Test suite failed:', error)
  process.exit(1)
})
