import 'package:flutter/material.dart';
import 'recipe.dart';

class RecipeCard extends StatefulWidget {
  final Recipe recipe;

  const RecipeCard({super.key, required this.recipe});

  @override
  _RecipeCardState createState() => _RecipeCardState();
}

class _RecipeCardState extends State<RecipeCard> {
  bool showIngredients = true;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        children: [
          ListTile(
            title: Text(widget.recipe.title),
            trailing: Switch(
              value: showIngredients,
              onChanged: (value) {
                setState(() {
                  showIngredients = value;
                });
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: showIngredients
                ? Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: widget.recipe.ingredients
                        .map((ingredient) => Text('- $ingredient'))
                        .toList(),
                  )
                : Text(widget.recipe.directions),
          ),
        ],
      ),
    );
  }
}
