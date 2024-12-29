import 'recipe.dart';

class ApiService {
  final String apiUrl = 'http://localhost:3000/recipes';

  Future<List<Recipe>> fetchRecipes() async {
    // final response = await http.get(Uri.parse(apiUrl));

    // if (response.statusCode == 200) {
    //   List<dynamic> data = json.decode(response.body);
    //   return data.map((json) => Recipe.fromJson(json)).toList();
    // } else {
    //   throw Exception('Failed to load recipes');
    // }

    // TODO metadata replace with api calls
    return [
      Recipe(
        id: '1',
        title: 'Spaghetti Carbonara',
        ingredients: [
          '200g spaghetti',
          '2 eggs',
          '50g pecorino cheese',
          '50g pancetta',
          '1 clove of garlic',
        ],
        directions:
            'Cook the spaghetti, mix the eggs and cheese, fry the pancetta and garlic, combine everything',
      ),
      Recipe(
        id: '2',
        title: 'Spaghetti Bolognese',
        ingredients: [
          '200g spaghetti',
          '200g minced beef',
          '1 onion',
          '1 clove of garlic',
          '1 can of tomatoes',
        ],
        directions:
            'Cook the spaghetti, fry the beef, onion, and garlic, add the tomatoes, combine everything',
      ),
    ];
  }
}
