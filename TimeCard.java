import java.io.File;
import java.io.FileNotFoundException;
import java.text.DecimalFormat;
import java.util.Scanner;

public class TimeCard {

    private static DecimalFormat moneyFormat = new DecimalFormat("##,###.##");
    private final short DATE = 0, DAY = 1, HOURS = 2, DRAWS = 3, NOTES = 4;    //constant indexes of each element of each line
    private final short MAX_COLUMNS = 5;
    private float totalHours, totalDraws, grossPay, wage;
    private String timeCardString, name, phone, email, start, pday;
    private short missedDays, period;


    public TimeCard(String csvPath, String startDate, String payDay, String payPeriodLength) throws FileNotFoundException
    {
        totalHours = 0;
        totalDraws = 0;
        missedDays = 0;
        grossPay = 0;
        wage = 0;
        timeCardString = "";
        name = "";
        phone = "";
        email = "";
        start = startDate;
        pday = payDay;
        period = Short.parseShort(payPeriodLength);
        if (period == 0)
            period = 7; //7 days is the default if 0 is given

        //Fill timecard with information from provided csv file
        buildTimeCard(csvPath);
    }

    private void buildTimeCard(String csvP) throws FileNotFoundException 
    {
        File csv = new File(csvP);
        Scanner scan = new Scanner(csv);
        String currentDate = "";
        String[] currentLine;

        //get necessary variables from first row of csv and assign them to instance variables
        currentLine = parseLine(scan.nextLine());
        name = currentLine[0];
        wage = Float.parseFloat(currentLine[1]);
        phone = currentLine[2];
        email = currentLine[3];


        //get headings for timeCardString
        timeCardString += buildTCHeader(); //top line with name, wage, payday
        currentLine = parseLine(scan.nextLine());
        timeCardString += buildTCLine(currentLine).replace('$', ' '); //second line with column headings

        //begin searching for startDate
        while (!currentDate.equals(start))
        {   
            //first line contains column headings so nothing gets missed on the first iteration
            if (scan.hasNextLine())
            {
                currentLine = parseLine(scan.nextLine());
                currentDate = currentLine[DATE];
            }
            else
                errorFound("Start date was never found. Check arguments.");
        }  

        int dayCounter = 0;  

        //startDate has been found, begin parsing and calculating
        for (int i = 0; i < period; i++) 
        {
            if (dayCounter > 16) //this is how I will catch some potential endDate error
                errorFound("Pay period was > 15 days. Check arguments.");
            else
                dayCounter++;
            
            timeCardString += buildTCLine(currentLine);

            if (!currentLine[DRAWS].equals(""))
                totalDraws += Float.parseFloat(currentLine[DRAWS]);

            if (!currentLine[HOURS].equals(""))
            {
                if (currentLine[HOURS].toLowerCase().contains("miss"))
                    missedDays++;
                else
                    totalHours += Float.parseFloat(currentLine[HOURS]);
            }

            currentLine = parseLine(scan.nextLine());
            currentDate = currentLine[DATE];
        }   

        scan.close();
        timeCardString += buildTCFooter();
    }

    // pull a line out of the csv and put each of its values into a string array
    private String[] parseLine(String line)
    {
        String[] result = new String[MAX_COLUMNS];
        Scanner lineScanner = new Scanner(line);
        lineScanner.useDelimiter(",");

        for (int i = 0; i < MAX_COLUMNS; i++)
            if (lineScanner.hasNext())
                result[i] = lineScanner.next();

        lineScanner.close();
        return result;
    }

    // build time card header
    private String buildTCHeader()
    {
        return phone + "\n" + email + "\n" + name + "\nPayday: " + pday + "\nWage: " + currency(moneyFormat.format(wage)) + "/hr\n\n";
    }

    // build time card body (line by line)
    private String buildTCLine(String[] line)
    {
        //make sure line[NOTES] isn't null
        String nts = line[NOTES] == null ? "" : line[NOTES];

        //make sure line[HOURS] doesn't have a rounding error
        //maximum possible hours for a day should be 5 characters -> ##.##
        String hrs = line[HOURS].length() > 5 ? line[HOURS].substring(0, 5) : line[HOURS];

        return String.format("%-9s| %-4s|%6s |%6s | %s\n", line[DATE], line[DAY], hrs, currency(line[DRAWS]), nts);
    }

    // build time card footer
    private String buildTCFooter()
    {
        String s = "\n";
        s += ("\nTotal Hours:\t" + moneyFormat.format(totalHours));
        s += ("\nTotal Draws:\t" + currency(moneyFormat.format(totalDraws)));

        if (missedDays > 0)
        {
            wage -= missedDays; //adjust wage according to # of missed days
            s += ("\nMissed Days:\t" + missedDays);
            s += ("\nAdjusted Wage:\t" + currency(moneyFormat.format(wage)) + "/hr");
        }

        grossPay = (totalHours * wage) - totalDraws;  //calculate gross pay
        s += ("\nGross Pay:\t" + currency(moneyFormat.format(grossPay)));

        return s;
    }
   

    private String currency(String num)
    {
        return num.equals("") ? num : "$" + num;
    }


    public String toString()
    {
        return timeCardString;
    }


    private void errorFound(String errorMessage)
    {
        System.err.println(errorMessage);
        System.exit(1);
    }
}
