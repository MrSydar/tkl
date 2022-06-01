# tkl

## Description
An application which allows you to upload your invoices to the `Księgowość360` platform ([link](https://www.360ksiegowosc.pl/)).

## How to use
![image](https://user-images.githubusercontent.com/50991602/171436556-3b40e1f2-ed1a-4f14-888a-074c85b164f5.png)

As an input, the application takes a CSV file with your invoices data. Select the file with `Select TKL report` button, click `Run` and wait until the upload process is finished.
When program finishes, there will be an `output.log` file created with logs so you can debug.

## CSV data

### Columns
1. `no`: invoice number
2. `date`: invoice date in `yyyyMMddHHmmss` format
3. `customer_nip`: NIP of the customer, leave empty if customer doesn't have it
4. `net`: net value
5. `tax`: tax value, so `net` + `tax` = gross
6. `tax_id`: tax id in the `Księgowość360` system
7. `customer_id`: customer id in the `Księgowość360` system
8. `product_code`: code of the product taken from the `Księgowość360` system
9. `product_description`: description the product. You can just copy it from the `Księgowość360` system

### Example
![image](https://user-images.githubusercontent.com/50991602/171439245-f2bd0205-23b6-448d-8865-faff0cd36e4c.png)
