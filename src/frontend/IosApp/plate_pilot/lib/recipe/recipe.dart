class Recipe {
  final String id;
  final String title;
  final List<String> ingredients;
  final String directions;

  Recipe({
    required this.id,
    required this.title,
    required this.ingredients,
    required this.directions,
  });

  factory Recipe.fromJson(Map<String, dynamic> json) {
    return Recipe(
      id: json['id'] as String,
      title: json['title'] as String,
      ingredients: List<String>.from(json['ingredients']),
      directions: json['directions'] as String,
    );
  }
}
