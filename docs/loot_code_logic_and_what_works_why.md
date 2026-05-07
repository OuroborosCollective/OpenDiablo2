# OpenDiablo2 Loot Code Logic, Mechanics, and What Works/Why

This document explores the deep code logic of the OpenDiablo2 loot and item drop systems, the mechanics behind monster drops and droplists, extra important nuances, and an overview of what currently works and what is still a work in progress (and why).

---

## 1. Droplists (Treasure Classes) Code Logic

In Diablo 2, "droplists" are managed by **Treasure Classes (TCs)**. In OpenDiablo2, these are represented by the `TreasureClassRecord` struct defined in `d2core/d2records/treasure_class_record.go`.

A single `TreasureClassRecord` consists of:
- `NumPicks`: The total number of items the engine attempts to generate from this TC.
  - If `NumPicks` is positive, the engine will roll a random treasure `NumPicks` times from the list of `Treasures`.
  - If `NumPicks` is negative, it guarantees counts for specific items instead of random probability rolls. The probability associated with each `Treasure` acts as a guaranteed drop count.
- `FreqNoDrop`: The weighted probability that *nothing* drops during a pick.
- `Treasures`: A slice of `Treasure` structs containing a `Code` (either an item or another TC) and a `Probability` weight.
- `FreqUnique`, `FreqSet`, `FreqRare`, `FreqMagic`: Base probabilities for determining item quality after it has been picked.

### Execution Flow: `ItemsFromTreasureClass`
When a monster dies or an object is opened, `ItemFactory.ItemsFromTreasureClass` (in `d2core/d2item/diablo2item/item_factory.go`) is invoked.
1. The engine evaluates `NumPicks`.
2. Based on whether `NumPicks` is positive or negative, it gathers a slice of `*d2records.Treasure`.
3. If `NumPicks` is positive, for each pick, it calls `rollTreasurePick()`, which calculates the total probability (including `FreqNoDrop`) and generates a random number.
4. If the random roll falls within an item's probability bucket, it adds that `Treasure` to a list of "Picks". If it hits `FreqNoDrop`, it moves to the next pick.
5. It recursively evaluates the gathered picks, turning `Treasure` codes into `Item`s.
6. As items bubble up from recursion, `applyDropModifier` is called to set the item's quality based on `rollDropModifier`.

---

## 2. Monster Drop Mechanics

Monsters use different Treasure Classes depending on their type, rank, and game difficulty.

### Standard Monsters (`MonsterStatsRecord`)
Regular monsters have different TC fields assigned in `MonStats.txt`, represented in `d2core/d2records/monster_stats_record.go`. Depending on difficulty and their rank, the engine pulls a different TC:
- **Normal**: `TreasureClassNormal`, `TreasureClassNightmare`, `TreasureClassHell`
- **Champions**: `TreasureClassChampionNormal`, `TreasureClassChampionNightmare`, `TreasureClassChampionHell`
- **Uniques**: `TreasureClass3UniqueNormal`, `TreasureClass3UniqueNightmare`, `TreasureClass3UniqueHell`
- **Quests**: `TreasureClassQuestNormal`, `TreasureClassQuestNightmare`, `TreasureClassQuestHell` (Requires specific Quest IDs to trigger).

### Super Unique Monsters (`SuperUniqueRecord`)
Super Uniques (e.g., Pindleskin, Corpsefire) have hardcoded mechanics defined in `d2core/d2records/monster_super_unique_record.go`. They bypass normal monster logic. They contain explicit `TreasureClassNormal`, `TreasureClassNightmare`, and `TreasureClassHell` fields. Super Uniques guarantee consistent, specific item drop logic distinct from standard unique monsters.

---

## 3. Extra Important Loot Mechanics

### Recursive Treasure Classes
Diablo 2's loot system is deeply recursive. When a `Treasure` is picked, OpenDiablo2 checks if the `Code` matches another Treasure Class.
- **Code Logic:** In `ItemsFromTreasureClass`, if the picked code exists in `f.asset.Records.Item.Treasure.Normal`, the engine recursively calls `ItemsFromTreasureClass(record)` on the subclass. This recursion continues until a base item (like `gld` or `armo`) is resolved. This architecture keeps data files clean; monsters simply reference an overarching TC (like "Act 1 Good"), which recursively drills down to specific items.

### Dynamic Item Level Resolution
Many TCs contain codes like `armo33` or `weap24` instead of specific items.
- **Code Logic:** `resolveDynamicTreasureCode` strips the numeric and alphabetic components. It then looks up equivalent item records (e.g., all items labeled `armo`) and selects one whose level falls within a 3-level range (e.g., levels 33, 34, 35). This allows a single TC entry to dynamically drop a variety of equivalently leveled items, saving massive amounts of data entry.

### Quality Modifiers and Fallbacks (Reconciliation)
After an item code is finalized, the engine determines its quality using `rollDropModifier()`.
- **Code Logic:** It maps out probabilities using a base of 1024, plus the `FreqUnique`, `FreqSet`, `FreqRare`, and `FreqMagic` from the TC.
- If an item rolls **Unique** but no Unique version exists for that base item, the system falls back to a **Rare** item with extra durability (or a Magic item if Rare fails). This graceful fallback matches the authentic Diablo 2 engine behavior and prevents the game from crashing or dropping an invalid item.
- If an item's underlying type record explicitly restricts its quality (e.g., items that can only be Normal or Magic), `sanitizeDropModifier` will downgrade the modifier appropriately.

---

## 4. What Works, What is WIP, and Why

### What Works:
1. **Recursive Droplist Architecture:** The engine successfully reads data dictionaries and accurately maps out all `TreasureClassRecords`. The recursive resolution of TCs all the way down to base items works precisely as the original game does.
2. **Negative `NumPicks` Logic:** The clever use of `NumPicks < 0` completely alters the loop behavior, turning probabilistic weights into guaranteed counts. This is critical for boss monsters (like the Countess) who are hardcoded to always drop a specific number of certain items (like runes or town portal scrolls).
3. **Item Factory & Drop Rolls:** The probability math, including `NoDrop` frequencies and `NumPicks`, accurately mirrors vanilla drop rates based solely on TC configurations.
4. **Dynamic Equivalency:** `armoXX` and `weapXX` mapping works perfectly, ensuring varied drops for armor and weapons.
5. **Base Quality Modifiers:** Generating Magic, Rare, Unique, and Set items based on TC probabilities is functional, and affixes are applied correctly to Magic and Rare items.
6. **Graceful Quality Fallbacks:** Generating Rare or Magic items when Unique or Set items fail works smoothly.

### What is WIP / Incomplete (And Why):
1. **Magic Find (MF):**
   - *Current State:* Not fully implemented in drop math.
   - *Why:* Currently, the drop modifier logic (`rollDropModifier`) relies *only* on the base probabilities defined in the TC. It does not yet inject the player's Magic Find stat or apply the necessary diminishing returns math required by Diablo 2.
2. **Item Ratio Records (`ItemRatio.txt`):**
   - *Current State:* The records are parsed into `ItemRatioRecord` (in `d2core/d2records/item_ratio_record.go`), but are not utilized heavily in the actual drop generation.
   - *Why:* In vanilla D2, item ratios dynamically alter base item drop rates taking into account the monster's level (mlvl), item level (ilvl), and MF. Integrating this requires tying the `ItemFactory` closer to the monster's active stats during combat, which is complex and still under development.
3. **Player Count (`players X` effect on NoDrop):**
   - *Current State:* `FreqNoDrop` is statically read from the TC.
   - *Why:* In vanilla, more players in a game exponentially reduce the `NoDrop` chance. This network/session state is not fully hooked into the item rolling logic yet.

## Summary
OpenDiablo2's loot engine successfully captures the foundational mechanics of Diablo 2—Treasure Class resolution, probability calculations, and dynamic item instantiation. The primary future work involves hooking external state factors (like Player Magic Find, Monster Level, and Player Count) into these foundational systems to achieve 1:1 drop accuracy.
