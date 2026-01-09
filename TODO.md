2. Multi user
    - first start with a single always logged in user so to speak? Or just integrate google oauth from the get go? Maybe email with keycloack for poc is easier? IDK
3, sync ios and frontend, or find a way to do e2e tests with ios app? 
4. come up with a epic roadmap, stuff is getting to cluttered
    - Recipes First! They are the backbone of the app
        - Especially the model should be refined. All fields that make up a recipe (regular data, metadata, usermetadata?)
        - creation and retrieval
        - editing can be done later
    - Reusable compontents to have unified style across app
        - have subnav items (like recipe cfreation and all other subnav items hide the stackbox)
    - Mealplan -> once core recipes is done and unified app
        - auto-suggest with constraints based on recipe properties
            - own metadata (how often suggested?)
        - save constraints for easy reuse
    - ShoppingList -> once core recipes is done and unified app
        - aggregate units and quantities
        - categories for ingredients -> should these live in the recipe data?
            - categories can be used for walking route
        - checkbox -> i prefer instant delete (maybe setting?)

<!-- once the above is done, we swarming -->
4. Features
    1. Agents
        - Chef agent which is responsible for combining meals, suggesting recipes
        - Fitness agent, has knowledge of selected meals and workout plans, can assist chef agent with selecting meals etc.
        - Unit & measurement prompt/agent for later shopping cart utils
