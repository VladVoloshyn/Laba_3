import java.util.Random;

public class HoneyPot {

    private static final int CAPACITY = 100;
    private static final int BEE_COUNT = 10;

    private static int honey = 0;
    private static int totalHoney = 0;
    private static Object lock = new Object();

    public static void main(String[] args) {
        Thread[] bees = new Thread[BEE_COUNT];


        for (int i = 0; i < BEE_COUNT; i++) {
            bees[i] = new Thread(new Bee());
            bees[i].start();
        }


        new Thread(new Bear()).start();
    }

    private static class Bee implements Runnable {

        @Override
        public void run() {
            while (totalHoney < CAPACITY) {

                int amount = new Random().nextInt(10);
                synchronized (lock) {
                    honey += amount;
                    totalHoney += amount;
                }


                if (honey == CAPACITY) {
                    synchronized (lock) {
                        lock.notify();
                    }
                }


                System.out.println("Бджола поклала " + amount + " ковтка меду");
            }
        }
    }

    private static class Bear implements Runnable {

        @Override
        public void run() {
            synchronized (lock) {
                while (totalHoney < CAPACITY) {
                    try {
                        lock.wait();
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }


                System.out.println("Ведмідь прокинувся");

                
                System.out.println("Ведмідь з'їв мед!");
                honey = 0;
            }
        }
    }
}
