#TODO combine all testing etc

run-backend:
	dotnet run --project src/backend/Hosting/Hosting.csproj
	
run-backend-watch:
	dotnet watch --project src/backend/Hosting/Hosting.csproj