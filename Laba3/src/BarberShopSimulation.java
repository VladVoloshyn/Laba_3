class BarberShop {
    private boolean customerWaiting = false;

    public synchronized void customerArrived() {
        customerWaiting = true;
        notify();
    }

    public synchronized void haircut() {
        while (!customerWaiting) {
            try {
                wait();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
        System.out.println("Перукар починає стрижку.");
        try {
            Thread.sleep(2000);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        System.out.println("Перукар завершив стрижку.");
        customerWaiting = false;
        notify();
    }
}

class Customer implements Runnable {
    private BarberShop barberShop;

    public Customer(BarberShop barberShop) {
        this.barberShop = barberShop;
    }

    @Override
    public void run() {
        System.out.println("Відвідувач прийшов в перукарню.");
        barberShop.customerArrived();
        barberShop.haircut();
        System.out.println("Відвідувач завершив стрижку та виходить.");
    }
}

public class BarberShopSimulation {
    public static void main(String[] args) {
        BarberShop barberShop = new BarberShop();


        Thread barberThread = new Thread(() -> {
            while (true) {
                barberShop.haircut();
            }
        });

        barberThread.start();


        for (int i = 1; i <= 5; i++)
        {
            Thread customerThread = new Thread(new Customer(barberShop));
            customerThread.start();
        }
    }
}
