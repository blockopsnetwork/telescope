--scrape-method) 
            SCRAPE_METHOD=$4
            if [[ $SCRAPE_METHOD="push" ]]; then
              REMOTE_WRITE_URL="https://telescope.blockops.network"
            elif [[ $SCRAPE_METHOD="pull" ]]; then
              if [[ -z $3 || $3 == -* ]]; then
                echo "Error: --scrape-method pull requires a REMOTE_WRITE_URL to be specified."
                exit 1
              else
                REMOTE_WRITE_URL=$3
                shift
              fi
            else
                echo "Error: Invalid scrape method. Please use 'push' or 'pull'."
                exit 1
            fi
            shift 2;;