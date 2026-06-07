package miao.byusi.pc.xylt;

import javafx.application.Application;
import javafx.fxml.FXMLLoader;
import javafx.scene.Parent;
import javafx.scene.Scene;
import javafx.stage.Stage;

public class MainApp extends Application {

    private static MainController mainController;
    private static Stage primaryStage;

    @Override
    public void start(Stage primaryStage) throws Exception {
        MainApp.primaryStage = primaryStage;

        FXMLLoader loader = new FXMLLoader(getClass().getResource("/fxml/main.fxml"));
        Parent root = loader.load();
        mainController = loader.getController();

        Scene scene = new Scene(root, 1200, 800);

        primaryStage.setTitle("校园论坛");
        primaryStage.setScene(scene);
        primaryStage.show();
    }

    public static MainController getMainController() {
        return mainController;
    }

    public static Stage getPrimaryStage() {
        return primaryStage;
    }

    public static void main(String[] args) {
        launch(args);
    }
}
