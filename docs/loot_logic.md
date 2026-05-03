# OpenDiablo2 Loot Logic and Drop Mechanics

This document provides an overview of how item drops, loot logic, and treasure classes work in the OpenDiablo2 engine based on the current codebase.

## 1. Treasure Classes (Droplists)

The core of the item drop system in OpenDiablo2 is the **Treasure Class (TC)**. TCs define a set of items or other TCs that can drop, along with the probabilities of each occurring. They are loaded from the game's data files (`TreasureClassEx.txt` / `TreasureClass.txt`) and represented in code by `TreasureClassRecord` (found in `d2core/d2records/treasure_class_record.go`).

A `TreasureClassRecord` contains:
- `NumPicks`: The number of items to attempt to drop from this class.
  - If `NumPicks` is positive, the engine will roll for an item `NumPicks` times.
  - If `NumPicks` is negative, each item's probability is treated as a guaranteed count for how many of that item will drop.
- `FreqNoDrop`: The probability that no item drops during a pick.
- `Treasures`: A list of potential drops (each being a `Treasure`), which can either be specific item codes (like `gld` for gold) or names of other Treasure Classes.
- **Drop Modifiers**: Probabilities for the dropped item to be of a specific quality (`FreqUnique`, `FreqSet`, `FreqRare`, `FreqMagic`).

### Recursive Drops
When a Treasure Class is rolled (`ItemsFromTreasureClass` in `d2core/d2item/diablo2item/item_factory.go`), it will pick a treasure based on the defined probabilities. If the picked treasure code corresponds to another Treasure Class, the engine recursively rolls from that sub-class until it resolves to a specific item or results in a `NoDrop`.

## 2. Monster Drop Mechanics

Monsters dictate which Treasure Class is used when they are killed. This is defined in their stat records.

### Standard Monsters (`MonsterStatsRecord`)
Regular monsters have different Treasure Classes depending on the difficulty and their rank:
- **Normal**: `TreasureClassNormal`, `TreasureClassNightmare`, `TreasureClassHell`
- **Champions**: `TreasureClassChampionNormal`, `TreasureClassChampionNightmare`, `TreasureClassChampionHell`
- **Uniques**: `TreasureClass3UniqueNormal`, `TreasureClass3UniqueNightmare`, `TreasureClass3UniqueHell`
- **Quests**: Special quest drops utilize `TreasureClassQuestNormal`, `TreasureClassQuestNightmare`, and `TreasureClassQuestHell` when specific quest trigger IDs are met.

### Super Unique Monsters (`SuperUniqueRecord`)
Super Unique monsters (like Pindleskin or Corpsefire) bypass the standard monster stats logic and have their own explicitly defined Treasure Classes for each difficulty (`TreasureClassNormal`, `TreasureClassNightmare`, `TreasureClassHell` in `d2core/d2records/monster_super_unique_record.go`).

## 3. Item Generation and Quality Logic

Once an item code is selected from a Treasure Class, the engine instantiates it via the `ItemFactory`.

### Dynamic Level Resolution
Some treasure codes are dynamic, like `armo33`. The engine parses these (via `resolveDynamicTreasureCode`) to select an equivalent item (e.g., armor) with a level within a range of 3 levels from the specified number (e.g., levels 33, 34, 35).

### Applying Drop Modifiers (Item Quality)
After an item is picked, the engine determines its quality (Unique, Set, Rare, Magic, or None/Normal). Currently, this roll is largely driven by the frequencies defined directly in the `TreasureClassRecord`.

The `rollDropModifier` function uses the following baseline for probability distribution:
1. `dropModifierBaseProbability` (which is 1024)
2. `FreqUnique`
3. `FreqSet`
4. `FreqRare`
5. `FreqMagic`

A random roll determines the quality. The `Item.applyDropModifier()` function then attempts to assign the appropriate affixes or unique/set properties.

### Reconciliation Fallbacks
If the engine rolls a specific quality but cannot find a valid record for it (e.g., it rolls Unique, but there is no Unique version of that base item), the engine employs reconciliation logic:
- A failed **Unique** or **Set** roll gracefully falls back to generating a **Rare** item (or a magic item with enhanced durability depending on limitations like `NoLimit`).
- If an item's underlying type record explicitly restricts its quality (e.g., items that can only be Normal or Magic), `sanitizeDropModifier` will downgrade the modifier appropriately to match the item type's allowed properties.

## 4. Item Ratio Records (Work in Progress)

The engine reads and parses `ItemRatio.txt` into `ItemRatioRecord` structs (`d2core/d2records/item_ratio_record.go`). In the original Diablo 2 logic, these ratios dynamically impact drop rates by taking into account the monster's level, player's Magic Find, and the item's base level.

Currently, the data dictionary loaders parse these into `ItemRatios` in the `RecordManager`, separating data for Normal, Exceptional/Uber, and Class-Specific items.
