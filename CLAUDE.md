# DarwinKit Development Guidelines

## Build Commands
- Test all packages: `go test ./...`
- Test a specific package: `go test ./macos/appkit`
- Generate bindings: `go generate ./...`
- Generate for a framework: `go generate ./macos/appkit`
- Run linting: `go vet ./...`
- Regenerate all frameworks: `./generate/tools/regen.sh macos`
- Clobber generated files: `go run ./generate/tools/clobbergen.go ./macos/[framework]`

## Code Style
- Use Go idiomatic code patterns when possible
- Follow standard Go formatting with `gofmt`
- Import ordering: standard library, then third-party, then internal packages
- Error handling: Always check errors returned by functions
- When generating bindings, ensure unique method signatures
- Follow existing naming conventions for generated code
- Memory management: Use `objc.Retain()` for objects that need to live beyond event loops
- Wrap code outside event loops in `objc.WithAutoreleasePool()`
- Follow binding API conventions (see docs/bindings.md):
  - Symbol prefixes are removed (NSWindow → Window)
  - Selector names converted to PascalCase (setFrame:display: → SetFrameDisplay)
  - Class methods get function variants with class name prefix

## Generation Pipeline
DarwinKit uses a multi-stage code generation pipeline to create Go bindings for Apple frameworks:

1. **symbolsdb**: The core database of Apple framework symbols (JSON files in a zip archive)
2. **declparse**: Parses Objective-C declarations into structured Go representations
3. **modules**: Defines framework metadata and dependency relationships
4. **typing**: Maps Objective-C types to Go types
5. **codegen**: Generates Go code from parsed declarations and type mappings
6. **enums**: Exports constants and enums from frameworks to Go code

Generation often requires an iterative approach to handle framework-specific quirks:
1. Generate → encounter panic for unknown situation
2. Update code to handle the situation (add to skip/abstract lists or fix parsing)
3. Generate again until successful

## Adding New Frameworks
Follow these detailed steps to add a new framework:

1. **Add to modules.go**:
   ```go
   {"FrameworkName", "Framework Display Name", "packagename", "FrameworkName/FrameworkName.h", []string{"NS", "Other prefixes"}},
   ```

2. **Handle Dependencies**:
   - Update coupling maps in modules.go if needed:
     - `CanAbstractModuleCoupling`: Makes types become generic interfaces
     - `CanSkipModuleCoupling`: Skips methods/properties using certain modules
     - `CanIgnoreNotFound`: List of modules that can be ignored if not found

3. **Export Constants**:
   ```
   go run ./generate/tools/enumexport.go [framework] > ./generate/modules/enums/macos/[framework]
   ```

4. **Initialize Framework Package**:
   ```
   go run ./generate/tools/initmod.go macos [framework]
   ```

5. **Generate Structs**:
   ```
   go run ./generate/tools/structs.go [framework] > ./macos/[framework]/[framework]_structs.go
   ```
   - Check the output file for any missing structs
   - Manually handle structs with `_Ctype_struct_` prefix by commenting out or replacing

6. **Generate Bindings**:
   ```
   go generate ./macos/[framework]
   ```
   - Handle panics by adding types to appropriate lists
   - Common issues include unknown types or dependency conflicts

7. **Test the Framework**:
   ```
   go test ./macos/[framework]
   ```
   - Fix compilation errors by handling missing types
   - Most often struct types from other frameworks

8. **Handle Circular Imports**:
   - Use `CanAbstractModuleCoupling` or `CanSkipModuleCoupling` to manage dependencies
   - For circular dependencies between frameworks, make one depend on the other as an abstract type

## Framework Coverage
- See API_ENHANCED_COVERAGE.md for detailed framework coverage analysis and statistics
- DarwinKit currently implements **41 frameworks** (approximately 45-50% of critical Apple frameworks)

### Implementation Status By Category
| Category | Implementation Level | Key Frameworks |
|----------|---------------------|----------------|
| Core | ~70% | Foundation (60%), CoreFoundation (95%) |
| UI | ~30% | AppKit (20%) |
| Graphics & Media | ~35% | CoreGraphics (50%), CoreImage, AVFoundation |
| Data & Storage | ~40% | CoreData, CoreSpotlight |
| Device Features | ~25% | CoreLocation, MapKit, HealthKit |
| Specialized APIs | ~15% | GameKit, JavaScriptCore, Security |

### Recently Added Frameworks (March 2025)
| Framework | Core Classes | Methods | Status |
|-----------|--------------|---------|--------|
| EventKit  | EventStore, Event, Calendar | 10+ | Basic Implementation |
| MapKit    | MapView, PointAnnotation | 10+ | Basic Implementation |
| HealthKit | HealthStore, QuantityType, Workout | 15+ | Basic Implementation |
| HomeKit   | HomeManager, Home, Room, Accessory | 20+ | Basic Implementation |
| GameKit   | Player, Leaderboard, Achievement | 25+ | Basic Implementation |
| Security  | Certificate, Identity, Trust, Policy | 15+ | Basic Implementation |
| JavaScriptCore | Context, Value, VirtualMachine | 25+ | Basic Implementation |
| UserNotifications | NotificationCenter, NotificationContent, NotificationTrigger | 20+ | Basic Implementation |

### Recommended Next Frameworks
1. **Core Services** - Critical for file system operations
2. **Metal** - Modern graphics API (high priority)
3. **CoreBluetooth** - Device connectivity
4. **NetworkExtension** - Advanced networking capabilities
5. **PDFKit** - PDF rendering and manipulation

## Tools Reference

### Generation Tools
- List framework constants: `go run ./generate/tools/constant.go macos [framework] [constant]`
- Verify declaration parsing: `go run ./generate/tools/declcheck.go [framework]`
- Check parsing coverage %: `go run ./generate/tools/declcheck.go [framework] | grep "Total coverage"`
- Check imports: `./generate/tools/imports.sh ./macos/[framework]`
- Search symbolsdb: `go run ./generate/tools/lookup.go [prefix]`
- Find symbol types: `go run ./generate/tools/type.go [typename]`
- View framework stats: `./generate/tools/stats.sh`
- Regenerate all frameworks: `./generate/tools/regen.sh macos`

### Analysis Tools
- Generate enhanced coverage report: `./cmd/tools/analyze_technologies.sh`
- Analyze Apple's frameworks: `go run cmd/tools/tech_analyzer/main.go --outdir="./analysis" --count-symbols --list-macos`
- Generate module entries: `go run cmd/tools/tech_analyzer/main.go --frameworks="[framework1],[framework2]" --gen-modules`
- Analyze implementation coverage: `go run cmd/tools/report_generator/main.go --coverage="$OUTPUT_DIR/analysis/coverage_report.json" --output="API_COVERAGE.md"`

## Implementation Strategy

When implementing a new framework or extending an existing one, follow these guidelines:

### Implementation Depth
- Focus on implementing the most commonly used classes and methods first
- Create "basic implementations" that include 10-25 core methods per class
- Prioritize symbols that are fundamental to the framework's functionality
- Implementation quality is more important than quantity

### Implementation Approach
1. **Staged Implementation**:
   - Start with core functionality that unlocks the framework's primary features
   - Add example code that demonstrates practical usage
   - Add more advanced features in subsequent iterations

2. **Dependency Handling**:
   - Favor abstract types over forcing implementation of dependent frameworks
   - Use CanAbstractModuleCoupling for optional dependencies
   - Use CanSkipModuleCoupling for rarely used dependencies

3. **Testing Strategy**:
   - Create simple tests that verify binding compilation
   - Create example applications that demonstrate real-world usage
   - Test integrations with other frameworks

## Documentation
- Add good comments for custom functionality
- Document memory management expectations for custom methods
- Reference Apple documentation URLs when appropriate
- Update API_ENHANCED_COVERAGE.md when adding new frameworks