@echo off
cd frontend
echo "Installing dependencies..."
npm install
echo "Building project..."
npm run build
echo "Project built successfully in the 'dist' folder."
