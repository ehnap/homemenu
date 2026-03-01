ALTER TABLE recipe_ingredients ADD COLUMN category TEXT DEFAULT 'ingredient';
ALTER TABLE recipes ADD COLUMN tips TEXT DEFAULT '';
