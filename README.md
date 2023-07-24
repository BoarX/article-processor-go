I had two approaches for this task One was to retrieve the list of articles periodicly then check if they exist in our database and pull the the articles that don't exist with their full data and store them in the database. however in this scenarion if we store an article in the database and for some reason it would get updated in the external provider, we would not get the update in our database. therefor I chose this solution.

TODO:
Add id logic for different clubs
Add checking for updated articles that are already stored in the database