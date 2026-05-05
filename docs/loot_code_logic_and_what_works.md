# Loot Code Logic and What Works

This document provides a deep dive into the code logic behind the item drop mechanics, droplists (Treasure Classes), and what makes the OpenDiablo2 loot system function correctly and efficiently.

## Loot Logic & Droplists (Treasure Classes)

The loot system heavily relies on "Treasure Classes" (TCs), which act as weighted droplists. The structure that defines these lists is `TreasureClassRecord`, found in `d2core/d2records/treasure_class_record.go`.

```go
type TreasureClassRecord struct {
	Name       string
	Group      int
	Level      int
	NumPicks   int
	FreqUnique int
	FreqSet    int
	FreqRare   int
	FreqMagic  int
	FreqNoDrop int
	Treasures  []*Treasure
}

type Treasure struct {
	Code        string
	Probability int
}
```

- **`NumPicks`**: Determines the logic used to pick items.
  - If `NumPicks` > 0: The engine will roll a random treasure `NumPicks` times from the list of `Treasures`.
  - If `NumPicks` < 0: The probability associated with each `Treasure` acts as a guaranteed drop count, dropping that specific item `Probability` times.
- **`FreqNoDrop`**: The weight/probability of nothing dropping during a single pick.
- **Quality Modifiers** (`FreqUnique`, `FreqSet`, `FreqRare`, `FreqMagic`): Weights used to decide the item's quality (e.g., Magic or Unique) after it has been picked.
- **`Treasures`**: A slice of `Treasure` structs. A `Treasure` code can refer either to a specific base item (e.g., `gld` for gold) or to *another* Treasure Class.

## Monster Drop Mechanics

When monsters are defeated, they drop items based on their associated Treasure Classes.

### Standard Monsters
The `MonsterStatsRecord` (`d2core/d2records/monster_stats_record.go`) holds references to Treasure Classes based on the monster's rank and the game difficulty:

- **Normal Monsters**: `TreasureClassNormal`, `TreasureClassNightmare`, `TreasureClassHell`
- **Champions**: `TreasureClassChampionNormal`, `TreasureClassChampionNightmare`, `TreasureClassChampionHell`
- **Uniques**: `TreasureClass3UniqueNormal`, `TreasureClass3UniqueNightmare`, `TreasureClass3UniqueHell`
- **Quests**: `TreasureClassQuestNormal`, `TreasureClassQuestNightmare`, `TreasureClassQuestHell` (Requires specific Quest IDs to trigger).

### Super Unique Monsters
Super Unique monsters (like Corpsefire or Pindleskin) have hardcoded TCs that bypass standard logic. These are defined in `SuperUniqueRecord` (`d2core/d2records/monster_super_unique_record.go`), explicitly containing `TreasureClassNormal`, `TreasureClassNightmare`, and `TreasureClassHell`.

## Code Logic Extraction

The core item generation happens within `ItemFactory` (`d2core/d2item/diablo2item/item_factory.go`).

### `ItemsFromTreasureClass`
This method recursively evaluates a `TreasureClassRecord`.
1. **Determine Picks**: Based on whether `NumPicks` is positive or negative, it gathers a slice of `*d2records.Treasure`.
2. **Recursive Resolution**: It iterates through the gathered picks:
   - If the `Code` matches an existing Treasure Class, it recursively calls `ItemsFromTreasureClass` to evaluate the sub-class.
   - If the `Code` matches an Item, it uses `ItemFromTreasure` to generate it.
3. **Apply Quality**: As items bubble up from the recursion, `applyDropModifier` is called to set the item's quality based on `rollDropModifier`.

### `rollTreasurePick`
This function performs the weighted random roll to pick a single `Treasure` from the `TreasureClassRecord`.
It builds an array of cumulative probabilities `tprob`, starting with `FreqNoDrop`, and then adds the `Probability` of each `Treasure`. A random number is rolled up to the `total`, and the code iterates through `tprob` to find the winning drop.

### `ItemFromTreasure`
Once a specific item `Treasure` code is selected, this function resolves it into an `Item`:
1. **Direct Match**: Checks if the code matches a common item exactly.
2. **Equivalency List**: Checks if the code is a generic category (like `armo`), returning a random item from that category's equivalency list.
3. **Dynamic Level Resolution**: Checks if the code is dynamic (e.g., `armo33`). It uses regex to split the string and numeric components and finds an equivalent item whose level falls within `numericComponent` and `numericComponent + 3`.

### `rollDropModifier`
Determines the quality of an item dropped from a Treasure Class.
It calculates cumulative probabilities based on a base probability (1024) and the `Freq` fields in `TreasureClassRecord`. A random number determines if the item becomes Unique, Set, Rare, Magic, or Normal.

## What Works & Why

1. **Recursive Droplist Architecture**
   - **Why it works**: Allowing a `Treasure` code to reference another `TreasureClassRecord` keeps data files like `TreasureClassEx.txt` clean and hierarchical. Instead of repeating common item pools (like "all level 15 weapons") across dozens of monsters, monsters simply reference an overarching "Act 1 Good" TC, which recursively drills down to specific items.
2. **Negative `NumPicks` Logic**
   - **Why it works**: The clever use of `NumPicks < 0` completely alters the loop behavior, turning probabilistic weights into guaranteed counts. This is critical for boss monsters (like the Countess) who are hardcoded to always drop a specific number of certain items (like runes or town portal scrolls).
3. **Dynamic Treasure Codes (e.g., `armo33`)**
   - **Why it works**: `ItemFromTreasure` dynamically parses codes to find items within a dynamic level range (`+3` levels). This saves massive amounts of data entry; rather than hardcoding every single armor piece a monster *could* drop into a TC, the system dynamically queries the asset manager for equivalent items that match the desired level band.
4. **Graceful Quality Fallbacks (Reconciliation)**
   - **Why it works**: The system rolls for quality *before* checking if the item can actually be that quality. If an item rolls "Unique" but no unique version of that base item exists, the `applyDropModifier` logic (along with `sanitizeDropModifier` found in `item.go`) gracefully downgrades the item to "Rare" (or "Magic" with increased durability), matching the authentic Diablo 2 engine behavior and preventing crashes or empty drops.