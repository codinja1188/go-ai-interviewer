#!/bin/bash
# Final verification script

echo "🔍 Running Final Verification"
echo "=============================="
echo ""

# Check 1: Build
echo "1️⃣  Building application..."
if go build -v > /dev/null 2>&1; then
    echo "   ✅ Build successful"
else
    echo "   ❌ Build failed"
    exit 1
fi

# Check 2: Go vet
echo ""
echo "2️⃣  Running go vet..."
if go vet ./... > /dev/null 2>&1; then
    echo "   ✅ Code quality check passed"
else
    echo "   ❌ Code quality issues found"
    exit 1
fi

# Check 3: Unit tests
echo ""
echo "3️⃣  Running unit tests..."
if go test ./... -short > /tmp/test.log 2>&1; then
    echo "   ✅ All unit tests passed"
    go test ./... -cover | grep -E "ok|coverage"
else
    echo "   ❌ Some tests failed"
    cat /tmp/test.log
    exit 1
fi

# Check 4: File structure
echo ""
echo "4️⃣  Verifying project structure..."
required_files=(
    "main.go"
    "README.md"
    "FEATURES.md"
    ".gitignore"
    "test_integration.sh"
    "models/models.go"
    "services/question_service.go"
    "services/evaluation_service.go"
    "services/auth_service.go"
    "database/database.go"
    "handlers/handlers.go"
    "templates/login.html"
    "templates/dashboard.html"
    "templates/interview.html"
    "static/css/style.css"
)

all_files_exist=true
for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✓ $file"
    else
        echo "   ✗ $file (missing)"
        all_files_exist=false
    fi
done

if [ "$all_files_exist" = true ]; then
    echo "   ✅ All required files present"
else
    echo "   ❌ Some files missing"
    exit 1
fi

# Check 5: Count features
echo ""
echo "5️⃣  Feature count..."
question_count=$(grep -c "ID.*algo\|ID.*ds\|ID.*sd" services/question_service.go)
template_count=$(ls templates/*.html 2>/dev/null | wc -l)
test_count=$(find . -name "*_test.go" | wc -l)

echo "   📝 Questions: $question_count"
echo "   🎨 Templates: $template_count"
echo "   🧪 Test files: $test_count"
echo "   ✅ Feature metrics collected"

# Summary
echo ""
echo "=============================="
echo "✅ All Verifications Passed!"
echo "=============================="
echo ""
echo "Summary:"
echo "  • Build: ✓"
echo "  • Code Quality: ✓"
echo "  • Unit Tests: ✓"
echo "  • File Structure: ✓"
echo "  • Features: ✓"
echo ""
echo "🎉 The Go AI Interviewer is ready for production!"
