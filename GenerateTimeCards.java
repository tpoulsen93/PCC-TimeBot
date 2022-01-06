import java.io.File;
import java.io.FileNotFoundException;
import java.text.ParseException;

public class GenerateTimeCards {

    private static String usage = "Usage: generateTimeCards <csv filepath> <start date> <pay date> [<number of days in pay period>] \nDate format: mm-dd-yy";

    public static void main(String[] args) throws ParseException, FileNotFoundException
    {
        //check arguments
        if (args.length < 4 || args.length > 4)
        {
            System.err.println(usage);
            System.exit(1);
        }

        String csvFileName = args[0];       //get csv file path
        String startDate = args[1];         //get startDate of pay period
        String payDay = args[2];            //get payDay date
        String payPeriodLength = args[3];   //get # of days in pay period

        TimeCard tc = new TimeCard(csvFileName, startDate, payDay, payPeriodLength);


        System.out.println(tc.toString());
    }
}
