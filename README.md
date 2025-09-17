# 🎲 Markovian Word Generator

A web application that generates new words based on character patterns from input vocabulary using Markov chains.

## 🚀 Quick Start

```bash
# Build and run
./build.sh
./build/markovian-word-gen

# Or run directly
go run .
```

Open http://localhost:8080 in your browser.

## 🧠 Algorithm Logic

### Sampler (`sampler.go`)

The `Sampler` is a generic weighted random sampling data structure that efficiently selects items based on their frequencies.

**Core Components:**
- `itemToIdxMap`: Maps items to their index in the items array
- `items`: Array storing all unique items
- `cumFreq`: Cumulative frequency array for efficient sampling

**How it works:**
1. **Adding items**: When adding an item with frequency `f`, it updates the cumulative frequency array
2. **Sampling**: Uses binary search on cumulative frequencies to select items proportionally to their weights
3. **Efficiency**: O(log n) sampling time with O(1) average case for updates

**Example:**
```go
sampler := NewSampler[rune](29)
sampler.Add('a', 3)  // 'a' appears 3 times
sampler.Add('b', 1)  // 'b' appears 1 time
// 'a' has 75% chance, 'b' has 25% chance
```

### WordGenerator (`wordgenerator.go`)

The `WordGenerator` uses Markov chains to generate words by learning character transition patterns.

**Core Logic:**
1. **Context Building**: For each word, it creates context windows of length 3 (configurable)
2. **Character Distribution**: For each context, it tracks which characters follow and their frequencies
3. **Word Generation**: Starts with '^' and generates characters based on learned patterns until '$'

**Process:**
```go
// Input: "apple" -> "^apple$"
// Contexts learned:
// "^" -> 'a' (1 time)
// "^a" -> 'p' (1 time)  
// "^ap" -> 'p' (1 time)
// "pp" -> 'l' (1 time)
// "pl" -> 'e' (1 time)
// "le" -> '$' (1 time)

// Generation: Start with "^", sample next character based on context
```

**Key Features:**
- **Context Fallback**: If a context doesn't exist, it shortens the context until one is found
- **Weighted Sampling**: Uses the Sampler for proportional character selection
- **Boundary Markers**: Uses '^' and '$' to mark word boundaries

## 📝 Source Code

- [WordGenerator Implementation](https://github.com/SaiNageswarS/markovian-word-gen/blob/master/wordgenerator.go)
- [Sampler Implementation](https://github.com/SaiNageswarS/markovian-word-gen/blob/master/sampler.go)